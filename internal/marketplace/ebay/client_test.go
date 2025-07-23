package ebay

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestEbayClient_Search(t *testing.T) {
    // Mock eBay API response
    mockResponse := `{
        "itemSummaries": [
            {
                "title": "Test Product 1",
                "price": {
                    "value": "29.99",
                    "currency": "USD"
                },
                "itemWebUrl": "https://ebay.com/item/123"
            },
            {
                "title": "Test Product 2",
                "price": {
                    "value": "49.99",
                    "currency": "USD"
                },
                "itemWebUrl": "https://ebay.com/item/456"
            }
        ]
    }`

    // Create test server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))
        assert.Contains(t, r.URL.Query().Get("q"), "laptop")
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(mockResponse))
    }))
    defer server.Close()

    // Create client with test server URL
    client := NewClient("test-api-key")
    client.baseURL = server.URL

    // Test search
    products, err := client.Search("laptop", 0, 100)

    assert.NoError(t, err)
    assert.Len(t, products, 2)
    assert.Equal(t, "Test Product 1", products[0].Title)
    assert.Equal(t, 29.99, products[0].Price)
    assert.Equal(t, "https://ebay.com/item/123", products[0].Link)
}

func TestEbayClient_Search_APIError(t *testing.T) {
    // Create test server that returns error
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusInternalServerError)
    }))
    defer server.Close()

    client := NewClient("test-api-key")
    client.baseURL = server.URL

    products, err := client.Search("laptop", 0, 100)

    assert.Error(t, err)
    assert.Nil(t, products)
    assert.Contains(t, err.Error(), "eBay API returned status 500")
}