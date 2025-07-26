package openai

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"
)

// handleSearchIntent processes search-only requests
func (c *Client) handleSearchIntent(ctx context.Context, intent MessageIntent) (*OrchestratorResponse, error) {
	if intent.Product == "" {
		return &OrchestratorResponse{
			Message: "I couldn't identify what product you're looking for. Please specify a product name.",
			Errors:  []string{"No product specified in search request"},
		}, nil
	}

	// Execute search with retry logic
	searchResult, err := c.executeWithRetry(ctx, FunctionCall{
		Name: "search_marketplace",
		Arguments: map[string]any{
			"query":     intent.Product,
			"min_price": intent.MinPrice,
			"max_price": intent.MaxPrice,
		},
	})
	if err != nil {
		return &OrchestratorResponse{
			Message: fmt.Sprintf("I encountered an error while searching for %s: %v", intent.Product, err),
			Errors:  []string{err.Error()},
		}, nil
	}

	// Convert search results
	searchSummary := c.convertSearchResults(searchResult, intent.Product)

	message := c.formatSearchMessage(searchSummary)

	return &OrchestratorResponse{
		Message:       message,
		SearchResults: searchSummary,
	}, nil
}

// handleMarketingIntent processes marketing-only requests
func (c *Client) handleMarketingIntent(ctx context.Context, intent MessageIntent) (*OrchestratorResponse, error) {
	if intent.Description == "" && intent.Product == "" {
		return &OrchestratorResponse{
			Message: "I need a product description to create marketing copy. Please provide more details about the product.",
			Errors:  []string{"No product description provided for marketing"},
		}, nil
	}

	description := intent.Description
	if description == "" {
		description = intent.Product
	}

	// Get taste profile
	tasteResult, err := c.executeWithRetry(ctx, FunctionCall{
		Name: "get_taste_profile",
		Arguments: map[string]any{
			"description": description,
		},
	})

	var segments []string
	if err != nil {
		log.Printf("Taste profile failed, using default segments: %v", err)
		segments = []string{"General Consumers", "Value Seekers"}
	} else {
		segments = c.extractSegments(tasteResult)
	}

	// Generate ad copy
	adResult, err := c.executeWithRetry(ctx, FunctionCall{
		Name: "generate_ad_copy",
		Arguments: map[string]any{
			"product_title": intent.Product,
			"segments":      segments,
		},
	})
	if err != nil {
		return &OrchestratorResponse{
			Message: fmt.Sprintf("I couldn't generate marketing copy: %v", err),
			Errors:  []string{err.Error()},
		}, nil
	}

	marketing := c.convertMarketingResults(adResult, segments)
	message := c.formatMarketingMessage(marketing, intent.Product)

	return &OrchestratorResponse{
		Message:   message,
		Marketing: marketing,
	}, nil
}

// handleCombinedIntent processes requests requiring both search and marketing
func (c *Client) handleCombinedIntent(ctx context.Context, intent MessageIntent) (*OrchestratorResponse, error) {
	response := &OrchestratorResponse{}
	var errors []string

	// Step 1: Search for products
	if intent.Product != "" {
		searchResult, err := c.executeWithRetry(ctx, FunctionCall{
			Name: "search_marketplace",
			Arguments: map[string]any{
				"query":     intent.Product,
				"min_price": intent.MinPrice,
				"max_price": intent.MaxPrice,
			},
		})

		if err != nil {
			errors = append(errors, fmt.Sprintf("Search failed: %v", err))
		} else {
			response.SearchResults = c.convertSearchResults(searchResult, intent.Product)
		}
	}

	// Step 2: Generate marketing copy
	description := intent.Description
	if description == "" {
		description = intent.Product
	}

	if description != "" {
		// Get taste profile
		tasteResult, err := c.executeWithRetry(ctx, FunctionCall{
			Name: "get_taste_profile",
			Arguments: map[string]any{
				"description": description,
			},
		})

		var segments []string
		if err != nil {
			log.Printf("Taste profile failed, using default segments: %v", err)
			segments = []string{"General Consumers", "Value Seekers"}
		} else {
			segments = c.extractSegments(tasteResult)
		}

		// Generate ad copy
		adResult, err := c.executeWithRetry(ctx, FunctionCall{
			Name: "generate_ad_copy",
			Arguments: map[string]any{
				"product_title": intent.Product,
				"segments":      segments,
			},
		})

		if err != nil {
			errors = append(errors, fmt.Sprintf("Marketing generation failed: %v", err))
		} else {
			response.Marketing = c.convertMarketingResults(adResult, segments)
		}
	}

	// Format combined message
	response.Message = c.formatCombinedMessage(response.SearchResults, response.Marketing, intent.Product)
	response.Errors = errors

	return response, nil
}

// handleUnknownIntent processes unclear requests using OpenAI
func (c *Client) handleUnknownIntent(ctx context.Context, message string) (*OrchestratorResponse, error) {
	// Use OpenAI to understand and respond to the message
	aiResponse, functionCalls, err := c.SendMessage(message)
	if err != nil {
		return &OrchestratorResponse{
			Message: "I'm sorry, I couldn't understand your request. Please try asking about product searches or marketing copy generation.",
			Errors:  []string{err.Error()},
		}, nil
	}

	response := &OrchestratorResponse{
		Message: aiResponse,
	}

	// Execute any function calls returned by OpenAI
	for _, fc := range functionCalls {
		result, err := c.executeWithRetry(ctx, fc)
		if err != nil {
			response.Errors = append(response.Errors, err.Error())
			continue
		}

		// Process results based on function type
		switch fc.Name {
		case "search_marketplace":
			response.SearchResults = c.convertSearchResults(result, "")
		case "generate_ad_copy":
			response.Marketing = c.convertMarketingResults(result, []string{})
		}
	}

	return response, nil
}

// executeWithRetry executes a function call with exponential backoff retry logic
func (c *Client) executeWithRetry(ctx context.Context, fc FunctionCall) (any, error) {
	const maxRetries = 3
	const baseDelay = 100 * time.Millisecond

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		result, err := c.ExecuteFunctionCall(fc)
		if err == nil {
			return result, nil
		}

		lastErr = err
		if attempt < maxRetries-1 {
			delay := time.Duration(float64(baseDelay) * math.Pow(2, float64(attempt)))
			log.Printf("Function call %s failed (attempt %d/%d), retrying in %v: %v",
				fc.Name, attempt+1, maxRetries, delay, err)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}

	return nil, fmt.Errorf("function call %s failed after %d attempts: %w", fc.Name, maxRetries, lastErr)
}
