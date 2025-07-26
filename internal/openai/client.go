package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jesee-kuya/blue/internal/marketplace"
	"github.com/jesee-kuya/blue/internal/marketplace/amazon"
	"github.com/jesee-kuya/blue/internal/marketplace/ebay"
	"github.com/jesee-kuya/blue/internal/qloo"
	"github.com/sashabaranov/go-openai"
)

// Client represents an OpenAI GPT-4o client with function calling capabilities
type Client struct {
	OpenaiClient *openai.Client
	Model        string
	AmazonClient marketplace.Client
	EbayClient   marketplace.Client
	QlooClient   *qloo.Client
	Timeout      time.Duration
}

// NewClient creates a new OpenAI client using the OPENAI_API_KEY environment variable
func NewClient() *Client {
	apiKey := os.Getenv("OPENAI_API_KEY")
	openaiClient := openai.NewClient(apiKey)

	return &Client{
		OpenaiClient: openaiClient,
		Model:        "gpt-4o",
		AmazonClient: amazon.NewClient("", "", "us-east-1"), // Mock credentials for now
		EbayClient:   ebay.NewClient(""),                    // Mock credentials for now
		QlooClient:   qloo.NewClient(),
		Timeout:      30 * time.Second,
	}
}

// NewClientWithKey creates a new OpenAI client with a specific API key (for testing)
func NewClientWithKey(apiKey string) *Client {
	openaiClient := openai.NewClient(apiKey)

	return &Client{
		OpenaiClient: openaiClient,
		Model:        "gpt-4o",
		AmazonClient: amazon.NewClient("", "", "us-east-1"), // Mock credentials for now
		EbayClient:   ebay.NewClient(""),                    // Mock credentials for now
		QlooClient:   qloo.NewClient(),
		Timeout:      30 * time.Second,
	}
}

// SendMessage sends a user message to GPT-4o and returns both text response and any function calls
func (c *Client) SendMessage(message string) (response string, functionCalls []FunctionCall, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	functions := GetFunctionDefinitions()
	tools := make([]openai.Tool, len(functions))
	for i, fn := range functions {
		tools[i] = openai.Tool{
			Type:     openai.ToolTypeFunction,
			Function: &fn,
		}
	}

	req := openai.ChatCompletionRequest{
		Model: c.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: message,
			},
		},
		Tools: tools,
	}

	resp, err := c.OpenaiClient.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", nil, fmt.Errorf("no response choices returned")
	}

	choice := resp.Choices[0]
	response = choice.Message.Content

	// Parse function calls if any
	if choice.Message.ToolCalls != nil {
		for _, toolCall := range choice.Message.ToolCalls {
			if toolCall.Type == openai.ToolTypeFunction {
				var args map[string]any
				if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
					return response, functionCalls, fmt.Errorf("failed to parse function arguments: %w", err)
				}

				functionCalls = append(functionCalls, FunctionCall{
					Name:      toolCall.Function.Name,
					Arguments: args,
				})
			}
		}
	}

	return response, functionCalls, nil
}

// ExecuteFunctionCall executes the requested function and returns results
func (c *Client) ExecuteFunctionCall(functionCall FunctionCall) (result any, err error) {
	switch functionCall.Name {
	case "search_marketplace":
		return c.executeSearchMarketplace(functionCall.Arguments)
	case "get_taste_profile":
		return c.executeGetTasteProfile(functionCall.Arguments)
	case "generate_ad_copy":
		return c.executeGenerateAdCopy(functionCall.Arguments)
	default:
		return nil, fmt.Errorf("unknown function: %s", functionCall.Name)
	}
}

// executeSearchMarketplace searches products across marketplaces
func (c *Client) executeSearchMarketplace(args map[string]any) (any, error) {
	var searchArgs SearchMarketplaceArgs

	query, ok := args["query"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid query parameter")
	}
	searchArgs.Query = query

	if minPrice, ok := args["min_price"].(float64); ok {
		searchArgs.MinPrice = minPrice
	}
	if maxPrice, ok := args["max_price"].(float64); ok {
		searchArgs.MaxPrice = maxPrice
	}

	// Search across multiple marketplaces
	var allProducts []marketplace.Product

	// Search Amazon
	amazonProducts, err := c.AmazonClient.Search(searchArgs.Query, searchArgs.MinPrice, searchArgs.MaxPrice)
	if err == nil {
		allProducts = append(allProducts, amazonProducts...)
	}

	// Search eBay
	ebayProducts, err := c.EbayClient.Search(searchArgs.Query, searchArgs.MinPrice, searchArgs.MaxPrice)
	if err == nil {
		allProducts = append(allProducts, ebayProducts...)
	}

	return map[string]any{
		"products": allProducts,
		"count":    len(allProducts),
	}, nil
}

// executeGetTasteProfile analyzes product description using Qloo API
func (c *Client) executeGetTasteProfile(args map[string]any) (any, error) {
	description, ok := args["description"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid description parameter")
	}

	segments, err := c.QlooClient.GetTasteProfile(description)
	if err != nil {
		return nil, fmt.Errorf("failed to get taste profile: %w", err)
	}

	return map[string]any{
		"segments": segments,
		"count":    len(segments),
	}, nil
}

// executeGenerateAdCopy generates marketing copy for target segments
func (c *Client) executeGenerateAdCopy(args map[string]any) (any, error) {
	productTitle, ok := args["product_title"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid product_title parameter")
	}

	segmentsInterface, ok := args["segments"].([]any)
	if !ok {
		return nil, fmt.Errorf("missing or invalid segments parameter")
	}

	segments := make([]string, len(segmentsInterface))
	for i, seg := range segmentsInterface {
		if segStr, ok := seg.(string); ok {
			segments[i] = segStr
		} else {
			return nil, fmt.Errorf("invalid segment type at index %d", i)
		}
	}

	// Generate ad copy using template-based approach
	adCopy := c.generateAdCopyTemplate(productTitle, segments)

	return adCopy, nil
}

// generateAdCopyTemplate creates ad copy using templates (simple implementation for now)
func (c *Client) generateAdCopyTemplate(productTitle string, segments []string) AdCopyResult {
	headlines := []string{
		fmt.Sprintf("Discover %s - Perfect for %s", productTitle, strings.Join(segments, " & ")),
		fmt.Sprintf("%s: Designed for %s", productTitle, segments[0]),
		fmt.Sprintf("Get Your %s Today!", productTitle),
	}

	descriptions := []string{
		fmt.Sprintf("Experience the best %s tailored for %s. Premium quality meets your unique needs.", productTitle, strings.Join(segments, ", ")),
		fmt.Sprintf("Join thousands of satisfied customers who chose %s. Perfect for %s looking for quality and value.", productTitle, segments[0]),
	}

	callToAction := "Shop Now and Transform Your Experience!"

	return AdCopyResult{
		Headlines:    headlines,
		Descriptions: descriptions,
		CallToAction: callToAction,
	}
}

// ProcessMessageWithTimeout processes a user message with orchestration and timeout
func (c *Client) ProcessMessageWithTimeout(message string, timeout time.Duration) (*OrchestratorResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return c.ProcessMessage(ctx, message)
}

// ProcessMessageSimple provides a simple interface for message processing
func (c *Client) ProcessMessageSimple(message string) (*OrchestratorResponse, error) {
	return c.ProcessMessageWithTimeout(message, c.Timeout)
}
