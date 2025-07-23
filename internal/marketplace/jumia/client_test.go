package jumia

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestJumiaClient_Search(t *testing.T) {
    // Mock Jumia HTML response
    mockHTML := `
        <html>
            <body>
                <h3 class="name">Test Laptop 1</h3>
                <div class="prc">₦ 150,000</div>
                <a href="/laptop-1" class="core">Link 1</a>
                
                <h3 class="name">Test Laptop 2</h3>
                <div class="prc">₦ 200,000</div>
                <a href="/laptop-2" class="core">Link 2</a>
            </body>
        </html>
    `

    // Create test server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        assert.Contains(t, r.URL.Query().Get("q"), "laptop")
        assert.Contains(t, r.Header.Get("User-Agent"), "Mozilla")
        w.Header().Set("Content-Type", "text/html")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(mockHTML))
    }))
    defer server.Close()

    // Create client with test server URL
    client := NewClient()
    client.baseURL = server.URL

    // Test search
    products, err := client.Search("laptop", 100000, 300000)

    assert.NoError(t, err)
    assert.Len(t, products, 2)
    assert.Equal(t, "Test Laptop 1", products[0].Title)
    assert.Equal(t, 150000.0, products[0].Price)
    assert.Equal(t, server.URL+"/laptop-1", products[0].Link)
}

func TestJumiaClient_Search_PriceFilter(t *testing.T) {
    mockHTML := `
        <html>
            <body>
                <h3 class="name">Cheap Laptop</h3>
                <div class="prc">₦ 50,000</div>
                <a href="/cheap-laptop" class="core">Link 1</a>
                
                <h3 class="name">Expensive Laptop</h3>
                <div class="prc">₦ 500,000</div>
                <a href="/expensive-laptop" class="core">Link 2</a>
            </body>
        </html>
    `

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(mockHTML))
    }))
    defer server.Close()

    client := NewClient()
    client.baseURL = server.URL

    // Test with price filter
    products, err := client.Search("laptop", 100000, 300000)

    assert.NoError(t, err)
    assert.Empty(t, products) // Both products should be filtered out
}

func TestJumiaClient_ParseProducts(t *testing.T) {
    client := NewClient()
    
    html := `
        <h3 class="name">Test Product</h3>
        <div class="prc">₦ 25,500</div>
        <a href="/test-product" class="core">Test Link</a>
    `

    products, err := client.parseProducts(html)

    assert.NoError(t, err)
    assert.Len(t, products, 1)
    assert.Equal(t, "Test Product", products[0].Title)
    assert.Equal(t, 25500.0, products[0].Price)
    assert.Equal(t, client.baseURL+"/test-product", products[0].Link)
}