package jumia

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jesee-kuya/blue/internal/cache"
	"github.com/jesee-kuya/blue/internal/marketplace"
)

// Client represents a Jumia API client
type Client struct {
	apiKey      string
	baseURL     string
	mockMode    bool
	redisClient *cache.RedisClient
}

// NewClient creates a new Jumia client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:      apiKey,
		baseURL:     "https://api.jumia.com/v1",
		mockMode:    true, // Mock mode for now
		redisClient: cache.NewRedisClient(),
	}
}

// Search searches for products on Jumia with Redis caching
func (c *Client) Search(query string, minPrice, maxPrice float64) ([]marketplace.Product, error) {
	// Generate cache key
	cacheKey := c.generateCacheKey(query, minPrice, maxPrice)
	
	// Try to get from cache first
	var cachedProducts []marketplace.Product
	if err := c.redisClient.Get(cacheKey, &cachedProducts); err == nil {
		return cachedProducts, nil
	}

	// Cache miss - fetch fresh data
	var products []marketplace.Product
	var err error
	
	if c.mockMode {
		products, err = c.mockSearch(query, minPrice, maxPrice)
	} else {
		return nil, fmt.Errorf("jumia API integration not yet implemented")
	}

	if err != nil {
		return nil, err
	}

	// Cache the results for 10 minutes
	c.redisClient.SetWithTTL(cacheKey, products, 10*time.Minute)
	
	return products, nil
}

// mockSearch provides mock data for testing and development
func (c *Client) mockSearch(query string, minPrice, maxPrice float64) ([]marketplace.Product, error) {
	// Simulate API delay
	time.Sleep(100 * time.Millisecond)

	// Generate mock products based on query
	mockProducts := []marketplace.Product{
		{
			Title: fmt.Sprintf("Jumia Deal: %s - Best Price", strings.Title(query)),
			Price: generatePrice(minPrice, maxPrice, 25.99),
			Link:  "https://jumia.com/product/12345",
		},
		{
			Title: fmt.Sprintf("Popular %s - Fast Delivery", strings.Title(query)),
			Price: generatePrice(minPrice, maxPrice, 35.99),
			Link:  "https://jumia.com/product/67890",
		},
	}

	// Filter by price range
	var filteredProducts []marketplace.Product
	for _, product := range mockProducts {
		if (minPrice == 0 || product.Price >= minPrice) &&
			(maxPrice == 0 || product.Price <= maxPrice) {
			filteredProducts = append(filteredProducts, product)
		}
	}

	return filteredProducts, nil
}

// generatePrice creates a realistic price within the given range
func generatePrice(minPrice, maxPrice, defaultPrice float64) float64 {
	if minPrice > 0 && maxPrice > 0 {
		return minPrice + rand.Float64()*(maxPrice-minPrice)
	}
	return defaultPrice + rand.Float64()*20 - 10 // ±10 variation
}

// generateCacheKey creates a cache key for the search parameters
func (c *Client) generateCacheKey(query string, minPrice, maxPrice float64) string {
	data := fmt.Sprintf("%s:%.2f:%.2f", query, minPrice, maxPrice)
	hash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
	return fmt.Sprintf("marketplace:search:jumia:%s:%.2f:%.2f", hash, minPrice, maxPrice)
}
