package qloo

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/jesee-kuya/blue/internal/cache"
)

// Client represents a Qloo Taste AI™ API client
type Client struct {
	apiKey      string
	baseURL     string
	httpClient  *http.Client
	redisClient *cache.RedisClient
}

// NewClient creates a new Qloo client using the QLOO_API_KEY environment variable
func NewClient() *Client {
	apiKey := os.Getenv("QLOO_API_KEY")
	return &Client{
		apiKey:      apiKey,
		baseURL:     "https://api.qloo.com/v1",
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		redisClient: cache.NewRedisClient(),
	}
}

// GetTasteProfile analyzes a product description with Redis caching
func (c *Client) GetTasteProfile(description string) ([]Segment, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("QLOO_API_KEY not set")
	}

	if description == "" {
		return []Segment{}, nil
	}

	// Generate cache key
	cacheKey := c.generateCacheKey(description)
	
	// Try to get from cache first
	var cachedSegments []Segment
	if err := c.redisClient.Get(cacheKey, &cachedSegments); err == nil {
		return cachedSegments, nil
	}

	// Cache miss - fetch fresh data
	segments, err := c.fetchTasteProfile(description)
	if err != nil {
		return nil, err
	}

	// Cache the results for 10 minutes
	c.redisClient.SetWithTTL(cacheKey, segments, 10*time.Minute)
	
	return segments, nil
}

// fetchTasteProfile performs the actual API call
func (c *Client) fetchTasteProfile(description string) ([]Segment, error) {
	request := TasteProfileRequest{Description: description}
	request.Options.MaxSegments = 10

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	reqURL := fmt.Sprintf("%s/taste/profile", c.baseURL)
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(payload))
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
		return nil, fmt.Errorf("qloo API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response TasteProfileResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return response.Segments, nil
}

// generateCacheKey creates a cache key for the description
func (c *Client) generateCacheKey(description string) string {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(description)))
	return fmt.Sprintf("qloo:profile:%s", hash)
}
