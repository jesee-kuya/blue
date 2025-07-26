package openai

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertSearchResults_EmptyResults(t *testing.T) {
	client := NewClientWithKey("test-key")

	result := map[string]any{
		"products": []any{},
		"count":    0,
	}

	summary := client.convertSearchResults(result, "test query")

	assert.Equal(t, "test query", summary.Query)
	assert.Equal(t, 0, summary.Count)
	assert.Empty(t, summary.Products)
}

func TestConvertSearchResults_InvalidFormat(t *testing.T) {
	client := NewClientWithKey("test-key")

	result := "invalid format"

	summary := client.convertSearchResults(result, "test query")

	assert.Equal(t, "test query", summary.Query)
	assert.Equal(t, 0, summary.Count)
	assert.Empty(t, summary.Products)
}

func TestExtractSegments_EmptyResult(t *testing.T) {
	client := NewClientWithKey("test-key")

	result := map[string]any{}

	segments := client.extractSegments(result)

	assert.Equal(t, []string{"General Consumers"}, segments)
}

func TestExtractSegments_InvalidFormat(t *testing.T) {
	client := NewClientWithKey("test-key")

	result := "invalid"

	segments := client.extractSegments(result)

	assert.Equal(t, []string{"General Consumers"}, segments)
}

func TestHandleSearchIntent_ExecutionError(t *testing.T) {
	mockClient := &MockClient{
		Client: NewClientWithKey("test-key"),
	}

	mockClient.mockExecuteFunc = func(fc FunctionCall) (any, error) {
		return nil, errors.New("API error")
	}

	intent := MessageIntent{
		Type:    IntentSearch,
		Product: "laptop",
	}

	ctx := context.Background()
	response, err := mockClient.handleSearchIntent(ctx, intent)

	assert.NoError(t, err) // Should not return error, but handle gracefully
	assert.NotNil(t, response)
	assert.Contains(t, response.Message, "error while searching")
	assert.Len(t, response.Errors, 1)
}

func TestHandleMarketingIntent_TasteProfileFailure(t *testing.T) {
	mockClient := &MockClient{
		Client: NewClientWithKey("test-key"),
	}

	callCount := 0
	mockClient.mockExecuteFunc = func(fc FunctionCall) (any, error) {
		callCount++
		switch fc.Name {
		case "get_taste_profile":
			return nil, errors.New("taste profile API error")
		case "generate_ad_copy":
			return AdCopyResult{
				Headlines:    []string{"Default Headline"},
				Descriptions: []string{"Default Description"},
				CallToAction: "Buy Now!",
			}, nil
		}
		return nil, nil
	}

	intent := MessageIntent{
		Type:        IntentMarketing,
		Product:     "Gaming Laptop",
		Description: "High-performance gaming laptop",
	}

	ctx := context.Background()
	response, err := mockClient.handleMarketingIntent(ctx, intent)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotNil(t, response.Marketing)
	assert.Equal(t, 2, callCount)                                        // Should still call ad copy generation
	assert.Contains(t, response.Marketing.Segments, "General Consumers") // Should use default segments
}

func TestHandleCombinedIntent_PartialFailure(t *testing.T) {
	mockClient := &MockClient{
		Client: NewClientWithKey("test-key"),
	}

	mockClient.mockExecuteFunc = func(fc FunctionCall) (any, error) {
		switch fc.Name {
		case "search_marketplace":
			return nil, errors.New("search failed")
		case "get_taste_profile":
			return map[string]any{
				"segments": []any{
					map[string]any{"name": "Tech Enthusiasts", "affinity_score": 0.85},
				},
			}, nil
		case "generate_ad_copy":
			return AdCopyResult{
				Headlines:    []string{"Ultimate Gaming Experience"},
				Descriptions: []string{"Experience gaming like never before"},
				CallToAction: "Shop Now!",
			}, nil
		}
		return nil, nil
	}

	intent := MessageIntent{
		Type:    IntentCombined,
		Product: "gaming laptop",
	}

	ctx := context.Background()
	response, err := mockClient.handleCombinedIntent(ctx, intent)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Nil(t, response.SearchResults) // Search failed
	assert.NotNil(t, response.Marketing)  // Marketing succeeded
	assert.Len(t, response.Errors, 1)     // Should have one error
	assert.Contains(t, response.Errors[0], "Search failed")
}

func TestFormatSearchMessage_NoResults(t *testing.T) {
	client := NewClientWithKey("test-key")

	results := &SearchResultsSummary{
		Query: "nonexistent product",
		Count: 0,
	}

	message := client.formatSearchMessage(results)

	assert.Contains(t, message, "couldn't find any products")
	assert.Contains(t, message, "nonexistent product")
}

func TestFormatSearchMessage_ManyResults(t *testing.T) {
	client := NewClientWithKey("test-key")

	products := make([]ProductSummary, 10)
	for i := 0; i < 10; i++ {
		products[i] = ProductSummary{
			Title: fmt.Sprintf("Product %d", i+1),
			Price: float64(100 + i*10),
			Link:  fmt.Sprintf("https://example.com/product%d", i+1),
		}
	}

	results := &SearchResultsSummary{
		Query:    "test products",
		Count:    10,
		Products: products,
	}

	message := client.formatSearchMessage(results)

	assert.Contains(t, message, "10 products")
	assert.Contains(t, message, "Product 1")
	assert.Contains(t, message, "Product 5")
	assert.Contains(t, message, "5 more results") // Should limit display
}

func TestIsStopWord(t *testing.T) {
	client := NewClientWithKey("test-key")

	stopWords := []string{"the", "and", "or", "in", "on", "at", "to", "for"}
	nonStopWords := []string{"laptop", "gaming", "wireless", "bluetooth"}

	for _, word := range stopWords {
		assert.True(t, client.isStopWord(word), "Expected '%s' to be a stop word", word)
	}

	for _, word := range nonStopWords {
		assert.False(t, client.isStopWord(word), "Expected '%s' to not be a stop word", word)
	}
}

func TestExtractProduct_ComplexMessage(t *testing.T) {
	client := NewClientWithKey("test-key")

	tests := []struct {
		message  string
		expected string
	}{
		{"Find me the best gaming laptops with RGB lighting", "best gaming laptops"},
		{"Search for wireless bluetooth headphones under $100", "wireless bluetooth headphones"},
		{"Show me professional cameras for photography", "professional cameras photography"},
		{"Get smartphones with good battery life", "smartphones good battery"},
	}

	for _, test := range tests {
		t.Run(test.message, func(t *testing.T) {
			product := client.extractProduct(test.message)
			assert.Contains(t, product, strings.Fields(test.expected)[0])
		})
	}
}
