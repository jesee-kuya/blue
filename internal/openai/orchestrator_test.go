package openai

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockClient extends Client for testing
type MockClient struct {
	*Client
	mockExecuteFunc func(FunctionCall) (any, error)
}

func TestClassifyIntent_SearchIntent(t *testing.T) {
	client := NewClientWithKey("test-key")

	tests := []struct {
		message  string
		expected IntentType
		product  string
	}{
		{"Find laptops under $1000", IntentSearch, "laptops"},
		{"Show me gaming headphones", IntentSearch, "gaming headphones"},
		{"Search for wireless mice between $20 and $50", IntentSearch, "wireless mice"},
		{"List smartphones", IntentSearch, "smartphones"},
	}

	for _, test := range tests {
		t.Run(test.message, func(t *testing.T) {
			intent := client.classifyIntent(test.message)
			assert.Equal(t, test.expected, intent.Type)
			assert.Contains(t, intent.Product, test.product)
		})
	}
}

func TestClassifyIntent_MarketingIntent(t *testing.T) {
	client := NewClientWithKey("test-key")

	tests := []struct {
		message  string
		expected IntentType
	}{
		{"Create marketing copy for gaming laptop", IntentMarketing},
		{"Generate ads for wireless headphones", IntentMarketing},
		{"Suggest marketing for smartphone", IntentMarketing},
		{"Make advertisement copy", IntentMarketing},
	}

	for _, test := range tests {
		t.Run(test.message, func(t *testing.T) {
			intent := client.classifyIntent(test.message)
			assert.Equal(t, test.expected, intent.Type)
		})
	}
}

func TestClassifyIntent_CombinedIntent(t *testing.T) {
	client := NewClientWithKey("test-key")

	tests := []struct {
		message  string
		expected IntentType
	}{
		{"Find laptops and create marketing copy", IntentCombined},
		{"Search for headphones and generate ads", IntentCombined},
		{"Show me phones and suggest marketing", IntentCombined},
	}

	for _, test := range tests {
		t.Run(test.message, func(t *testing.T) {
			intent := client.classifyIntent(test.message)
			assert.Equal(t, test.expected, intent.Type)
		})
	}
}

func TestExtractPriceRange(t *testing.T) {
	client := NewClientWithKey("test-key")

	tests := []struct {
		message  string
		minPrice float64
		maxPrice float64
	}{
		{"Find laptops under $1000", 0, 1000},
		{"Show phones less than $500", 0, 500},
		{"Search between $100 and $300", 100, 300},
		{"Items between $50 to $150", 50, 150},
		{"No price mentioned", 0, 0},
	}

	for _, test := range tests {
		t.Run(test.message, func(t *testing.T) {
			minPrice, maxPrice := client.extractPriceRange(test.message)
			assert.Equal(t, test.minPrice, minPrice)
			assert.Equal(t, test.maxPrice, maxPrice)
		})
	}
}

func TestHandleSearchIntent_Success(t *testing.T) {
	mockClient := &MockClient{
		Client: NewClientWithKey("test-key"),
	}

	mockClient.mockExecuteFunc = func(fc FunctionCall) (any, error) {
		if fc.Name == "search_marketplace" {
			return map[string]any{
				"products": []any{
					map[string]any{
						"title": "Gaming Laptop Pro",
						"price": 999.99,
						"link":  "https://example.com/laptop1",
					},
					map[string]any{
						"title": "Budget Laptop",
						"price": 599.99,
						"link":  "https://example.com/laptop2",
					},
				},
				"count": 2,
			}, nil
		}
		return nil, nil
	}

	intent := MessageIntent{
		Type:     IntentSearch,
		Product:  "laptop",
		MaxPrice: 1000,
	}

	ctx := context.Background()
	response, err := mockClient.handleSearchIntent(ctx, intent)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotNil(t, response.SearchResults)
	assert.Equal(t, 2, response.SearchResults.Count)
	assert.Contains(t, response.Message, "2 products")
	assert.Contains(t, response.Message, "laptop")
}

func TestHandleSearchIntent_NoProduct(t *testing.T) {
	client := NewClientWithKey("test-key")

	intent := MessageIntent{
		Type: IntentSearch,
	}

	ctx := context.Background()
	response, err := client.handleSearchIntent(ctx, intent)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Contains(t, response.Message, "couldn't identify what product")
	assert.Len(t, response.Errors, 1)
}

func TestHandleMarketingIntent_Success(t *testing.T) {
	mockClient := &MockClient{
		Client: NewClientWithKey("test-key"),
	}

	callCount := 0
	mockClient.mockExecuteFunc = func(fc FunctionCall) (any, error) {
		callCount++
		switch fc.Name {
		case "get_taste_profile":
			return map[string]any{
				"segments": []any{
					map[string]any{"name": "Tech Enthusiasts", "affinity_score": 0.85},
					map[string]any{"name": "Gamers", "affinity_score": 0.78},
				},
			}, nil
		case "generate_ad_copy":
			return AdCopyResult{
				Headlines:    []string{"Ultimate Gaming Experience", "Power Meets Performance"},
				Descriptions: []string{"Experience gaming like never before"},
				CallToAction: "Shop Now!",
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
	assert.Equal(t, 2, callCount) // Should call both taste profile and ad copy
	assert.Len(t, response.Marketing.Headlines, 2)
	assert.Contains(t, response.Message, "Gaming Laptop")
}

func TestHandleCombinedIntent_Success(t *testing.T) {
	mockClient := &MockClient{
		Client: NewClientWithKey("test-key"),
	}

	mockClient.mockExecuteFunc = func(fc FunctionCall) (any, error) {
		switch fc.Name {
		case "search_marketplace":
			return map[string]any{
				"products": []any{
					map[string]any{
						"title": "Gaming Laptop Pro",
						"price": 999.99,
						"link":  "https://example.com/laptop1",
					},
				},
				"count": 1,
			}, nil
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
	assert.NotNil(t, response.SearchResults)
	assert.NotNil(t, response.Marketing)
	assert.Contains(t, response.Message, "Product Listings")
	assert.Contains(t, response.Message, "Marketing Copy")
}

func TestExecuteWithRetry_Success(t *testing.T) {
	mockClient := &MockClient{
		Client: NewClientWithKey("test-key"),
	}

	callCount := 0
	mockClient.mockExecuteFunc = func(fc FunctionCall) (any, error) {
		callCount++
		return "success", nil
	}

	ctx := context.Background()
	fc := FunctionCall{Name: "test_function"}

	result, err := mockClient.executeWithRetry(ctx, fc)

	assert.NoError(t, err)
	assert.Equal(t, "success", result)
	assert.Equal(t, 1, callCount)
}

func TestExecuteWithRetry_SuccessAfterRetry(t *testing.T) {
	mockClient := &MockClient{
		Client: NewClientWithKey("test-key"),
	}

	callCount := 0
	mockClient.mockExecuteFunc = func(fc FunctionCall) (any, error) {
		callCount++
		if callCount < 2 {
			return nil, assert.AnError
		}
		return "success", nil
	}

	ctx := context.Background()
	fc := FunctionCall{Name: "test_function"}

	result, err := mockClient.executeWithRetry(ctx, fc)

	assert.NoError(t, err)
	assert.Equal(t, "success", result)
	assert.Equal(t, 2, callCount)
}

func TestExecuteWithRetry_MaxRetriesExceeded(t *testing.T) {
	mockClient := &MockClient{
		Client: NewClientWithKey("test-key"),
	}

	callCount := 0
	mockClient.mockExecuteFunc = func(fc FunctionCall) (any, error) {
		callCount++
		return nil, assert.AnError
	}

	ctx := context.Background()
	fc := FunctionCall{Name: "test_function"}

	result, err := mockClient.executeWithRetry(ctx, fc)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 3, callCount) // Should retry 3 times
	assert.Contains(t, err.Error(), "failed after 3 attempts")
}

func TestExecuteWithRetry_ContextCancellation(t *testing.T) {
	mockClient := &MockClient{
		Client: NewClientWithKey("test-key"),
	}

	mockClient.mockExecuteFunc = func(fc FunctionCall) (any, error) {
		time.Sleep(100 * time.Millisecond)
		return nil, assert.AnError
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	fc := FunctionCall{Name: "test_function"}

	result, err := mockClient.executeWithRetry(ctx, fc)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, context.DeadlineExceeded, err)
}

func TestProcessMessage_Integration(t *testing.T) {
	mockClient := &MockClient{
		Client: NewClientWithKey("test-key"),
	}

	mockClient.mockExecuteFunc = func(fc FunctionCall) (any, error) {
		switch fc.Name {
		case "search_marketplace":
			return map[string]any{
				"products": []any{
					map[string]any{
						"title": "Gaming Laptop Pro",
						"price": 999.99,
						"link":  "https://example.com/laptop1",
					},
				},
				"count": 1,
			}, nil
		}
		return nil, nil
	}

	ctx := context.Background()
	response, err := mockClient.ProcessMessage(ctx, "Find gaming laptops under $1000")

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotNil(t, response.SearchResults)
	assert.Equal(t, 1, response.SearchResults.Count)
	assert.Contains(t, response.Message, "gaming laptops")
}
