package amazon

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestAmazonClient_Search_MockMode(t *testing.T) {
    client := NewClient("test-access-key", "test-secret-key", "us-east-1")

    products, err := client.Search("laptop", 20, 50)

    assert.NoError(t, err)
    assert.NotEmpty(t, products)
    
    // Check that all products are within price range
    for _, product := range products {
        assert.GreaterOrEqual(t, product.Price, 20.0)
        assert.LessOrEqual(t, product.Price, 50.0)
        assert.NotEmpty(t, product.Title)
        assert.NotEmpty(t, product.Link)
        assert.Contains(t, product.Title, "Laptop")
    }
}

func TestAmazonClient_Search_NoPriceFilter(t *testing.T) {
    client := NewClient("test-access-key", "test-secret-key", "us-east-1")

    products, err := client.Search("phone", 0, 0)

    assert.NoError(t, err)
    assert.NotEmpty(t, products)
    
    for _, product := range products {
        assert.NotEmpty(t, product.Title)
        assert.Greater(t, product.Price, 0.0)
        assert.NotEmpty(t, product.Link)
        assert.Contains(t, product.Title, "Phone")
    }
}

func TestGeneratePrice(t *testing.T) {
    // Test with price range
    price := generatePrice(10, 20, 15)
    assert.GreaterOrEqual(t, price, 10.0)
    assert.LessOrEqual(t, price, 20.0)

    // Test with no range (should use default with variation)
    price = generatePrice(0, 0, 25)
    assert.Greater(t, price, 0.0)
}