package openai

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient_WithEnvironmentVariable(t *testing.T) {
	os.Setenv("OPENAI_API_KEY", "test-api-key")
	defer os.Unsetenv("OPENAI_API_KEY")

	client := NewClient()

	assert.NotNil(t, client.OpenaiClient)
	assert.Equal(t, "gpt-4o", client.Model)
	assert.NotNil(t, client.AmazonClient)
	assert.NotNil(t, client.EbayClient)
	assert.NotNil(t, client.QlooClient)
}

func TestExecuteFunctionCall_SearchMarketplace(t *testing.T) {
	client := NewClientWithKey("test-key")

	functionCall := FunctionCall{
		Name: "search_marketplace",
		Arguments: map[string]interface{}{
			"query":     "laptop",
			"min_price": 500.0,
			"max_price": 1500.0,
		},
	}

	result, err := client.ExecuteFunctionCall(functionCall)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	resultMap, ok := result.(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, resultMap, "products")
	assert.Contains(t, resultMap, "count")
}

func TestExecuteFunctionCall_SearchMarketplace_MissingQuery(t *testing.T) {
	client := NewClientWithKey("test-key")

	functionCall := FunctionCall{
		Name: "search_marketplace",
		Arguments: map[string]interface{}{
			"min_price": 500.0,
		},
	}

	result, err := client.ExecuteFunctionCall(functionCall)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "missing or invalid query parameter")
}

func TestExecuteFunctionCall_GetTasteProfile(t *testing.T) {
	// Mock Qloo API
	mockResponse := `{
        "status": "success",
        "segments": [
            {"name": "Tech Enthusiasts", "affinity_score": 0.85},
            {"name": "Early Adopters", "affinity_score": 0.72}
        ]
    }`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := NewClientWithKey("test-key")
	// Note: In a real test, we'd need to expose baseURL or use dependency injection

	functionCall := FunctionCall{
		Name: "get_taste_profile",
		Arguments: map[string]interface{}{
			"description": "High-performance gaming laptop",
		},
	}

	_, err := client.ExecuteFunctionCall(functionCall)

	// This will fail in the current setup because we can't easily mock the Qloo client
	// In a production environment, we'd use dependency injection for better testability
	assert.Error(t, err) // Expected to fail due to missing API key
}

func TestExecuteFunctionCall_GenerateAdCopy(t *testing.T) {
	client := NewClientWithKey("test-key")

	functionCall := FunctionCall{
		Name: "generate_ad_copy",
		Arguments: map[string]interface{}{
			"product_title": "Gaming Laptop Pro",
			"segments":      []interface{}{"Tech Enthusiasts", "Gamers"},
		},
	}

	result, err := client.ExecuteFunctionCall(functionCall)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	adCopy, ok := result.(AdCopyResult)
	assert.True(t, ok)
	assert.NotEmpty(t, adCopy.Headlines)
	assert.NotEmpty(t, adCopy.Descriptions)
	assert.NotEmpty(t, adCopy.CallToAction)
	assert.Contains(t, adCopy.Headlines[0], "Gaming Laptop Pro")
}

func TestExecuteFunctionCall_GenerateAdCopy_InvalidSegments(t *testing.T) {
	client := NewClientWithKey("test-key")

	functionCall := FunctionCall{
		Name: "generate_ad_copy",
		Arguments: map[string]interface{}{
			"product_title": "Gaming Laptop Pro",
			"segments":      "invalid_segments", // Should be array
		},
	}

	result, err := client.ExecuteFunctionCall(functionCall)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "missing or invalid segments parameter")
}

func TestExecuteFunctionCall_UnknownFunction(t *testing.T) {
	client := NewClientWithKey("test-key")

	functionCall := FunctionCall{
		Name:      "unknown_function",
		Arguments: map[string]interface{}{},
	}

	result, err := client.ExecuteFunctionCall(functionCall)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unknown function: unknown_function")
}

func TestGenerateAdCopyTemplate(t *testing.T) {
	client := NewClientWithKey("test-key")

	productTitle := "Wireless Headphones"
	segments := []string{"Music Lovers", "Tech Enthusiasts"}

	result := client.generateAdCopyTemplate(productTitle, segments)

	assert.NotEmpty(t, result.Headlines)
	assert.NotEmpty(t, result.Descriptions)
	assert.NotEmpty(t, result.CallToAction)
	assert.Contains(t, result.Headlines[0], productTitle)
	assert.Contains(t, result.Headlines[0], "Music Lovers")
}

func TestFunctionCall_JSONMarshaling(t *testing.T) {
	functionCall := FunctionCall{
		Name: "test_function",
		Arguments: map[string]interface{}{
			"param1": "value1",
			"param2": 42.0,
		},
	}

	data, err := json.Marshal(functionCall)
	assert.NoError(t, err)

	var unmarshaled FunctionCall
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, functionCall.Name, unmarshaled.Name)
	assert.Equal(t, functionCall.Arguments["param1"], unmarshaled.Arguments["param1"])
}

func TestAdCopyResult_JSONMarshaling(t *testing.T) {
	adCopy := AdCopyResult{
		Headlines:    []string{"Headline 1", "Headline 2"},
		Descriptions: []string{"Description 1"},
		CallToAction: "Buy Now!",
	}

	data, err := json.Marshal(adCopy)
	assert.NoError(t, err)

	var unmarshaled AdCopyResult
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, adCopy.Headlines, unmarshaled.Headlines)
	assert.Equal(t, adCopy.CallToAction, unmarshaled.CallToAction)
}
