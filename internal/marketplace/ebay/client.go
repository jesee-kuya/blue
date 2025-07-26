package ebay

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jesee-kuya/blue/internal/cache"
	"github.com/jesee-kuya/blue/internal/marketplace"
)

// Client represents an eBay API client
type Client struct {
	apiKey      string
	baseURL     string
	httpClient  *http.Client
	redisClient *cache.RedisClient
}

// NewClient creates a new eBay client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:      apiKey,
		baseURL:     "https://api.ebay.com/buy/browse/v1",
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		redisClient: cache.NewRedisClient(),
	}
}

// eBayResponse represents the response structure from eBay Browse API
type eBayResponse struct {
	ItemSummaries []struct {
		Title string `json:"title"`
		Price struct {
			Value    string `json:"value"`
			Currency string `json:"currency"`
		} `json:"price"`
		ItemWebURL string `json:"itemWebUrl"`
	} `json:"itemSummaries"`
}

// Search searches for products on eBay with Redis caching
func (c *Client) Search(query string, minPrice, maxPrice float64) ([]marketplace.Product, error) {
	// Generate cache key
	cacheKey := c.generateCacheKey(query, minPrice, maxPrice)

	// Try to get from cache first
	var cachedProducts []marketplace.Product
	if err := c.redisClient.Get(cacheKey, &cachedProducts); err == nil {
		return cachedProducts, nil
	}

	// Cache miss - fetch fresh data
	products, err := c.searchAPI(query, minPrice, maxPrice)
	if err != nil {
		return nil, err
	}

	// Cache the results for 10 minutes
	c.redisClient.SetWithTTL(cacheKey, products, 10*time.Minute)

	return products, nil
}

// searchAPI performs the actual API call
func (c *Client) searchAPI(query string, minPrice, maxPrice float64) ([]marketplace.Product, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("q", query)
	params.Set("limit", "50")

	if minPrice > 0 {
		params.Set("filter", fmt.Sprintf("price:[%s..%s],priceCurrency:USD",
			strconv.FormatFloat(minPrice, 'f', 2, 64),
			strconv.FormatFloat(maxPrice, 'f', 2, 64)))
	}

	// Create and execute request
	reqURL := fmt.Sprintf("%s/item_summary/search?%s", c.baseURL, params.Encode())
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("eBay API returned status %d", resp.StatusCode)
	}

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var ebayResp eBayResponse
	if err := json.Unmarshal(body, &ebayResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to standard format
	products := make([]marketplace.Product, 0, len(ebayResp.ItemSummaries))
	for _, item := range ebayResp.ItemSummaries {
		price, err := strconv.ParseFloat(item.Price.Value, 64)
		if err != nil {
			continue
		}

		products = append(products, marketplace.Product{
			Title: item.Title,
			Price: price,
			Link:  item.ItemWebURL,
		})
	}

	return products, nil
}

// generateCacheKey creates a cache key for the search parameters
func (c *Client) generateCacheKey(query string, minPrice, maxPrice float64) string {
	data := fmt.Sprintf("%s:%.2f:%.2f", query, minPrice, maxPrice)
	hash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
	return fmt.Sprintf("marketplace:search:ebay:%s:%.2f:%.2f", hash, minPrice, maxPrice)
}
