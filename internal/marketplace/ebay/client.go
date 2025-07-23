package ebay

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jesee-kuya/blue/internal/marketplace"
)

// Client represents an eBay API client
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new eBay client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "https://api.ebay.com/buy/browse/v1",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
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

// Search searches for products on eBay using the Browse API
func (c *Client) Search(query string, minPrice, maxPrice float64) ([]marketplace.Product, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("q", query)
	params.Set("limit", "50")

	if minPrice > 0 {
		params.Set("filter", fmt.Sprintf("price:[%s..%s],priceCurrency:USD",
			strconv.FormatFloat(minPrice, 'f', 2, 64),
			strconv.FormatFloat(maxPrice, 'f', 2, 64)))
	}

	// Create request
	reqURL := fmt.Sprintf("%s/item_summary/search?%s", c.baseURL, params.Encode())
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make request
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
			continue // Skip items with invalid prices
		}

		products = append(products, marketplace.Product{
			Title: item.Title,
			Price: price,
			Link:  item.ItemWebURL,
		})
	}

	return products, nil
}
