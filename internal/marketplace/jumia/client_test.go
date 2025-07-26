package jumia

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJumiaClient_Search_PriceFilter(t *testing.T) {
	mockHTML := `
        <html>
            <body>
                <h3 class="name">Cheap Laptop</h3>
                <div class="prc">₦ 50,000</div
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

	client := NewClient("test-api-key")
	client.baseURL = server.URL

	// Test with price filter
	products, err := client.Search("laptop", 100000, 300000)

	assert.NoError(t, err)
	assert.Empty(t, products) // Both products should be filtered out
}

