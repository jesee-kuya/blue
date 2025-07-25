package qloo

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "time"
)

// Client represents a Qloo Taste AI™ API client
type Client struct {
    apiKey     string
    baseURL    string
    httpClient *http.Client
}

// NewClient creates a new Qloo client using the QLOO_API_KEY environment variable
func NewClient() *Client {
    apiKey := os.Getenv("QLOO_API_KEY")
    return &Client{
        apiKey:  apiKey,
        baseURL: "https://api.qloo.com/v1",
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

// NewClientWithKey creates a new Qloo client with a specific API key
func NewClientWithKey(apiKey string) *Client {
    return &Client{
        apiKey:  apiKey,
        baseURL: "https://api.qloo.com/v1",
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

// GetTasteProfile analyzes a product description and returns audience segments with affinity scores
func (c *Client) GetTasteProfile(description string) ([]Segment, error) {
    if c.apiKey == "" {
        return nil, fmt.Errorf("QLOO_API_KEY not set")
    }

    if description == "" {
        return []Segment{}, nil
    }

    // Prepare request payload
    request := TasteProfileRequest{
        Description: description,
    }
    request.Options.MaxSegments = 10

    payload, err := json.Marshal(request)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }

    // Create HTTP request
    reqURL := fmt.Sprintf("%s/taste/profile", c.baseURL)
    req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(payload))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    // Set headers
    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")

    // Make request
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to make request: %w", err)
    }
    defer resp.Body.Close()

    // Read response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }

    // Handle non-200 status codes
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Qloo API returned status %d: %s", resp.StatusCode, string(body))
    }

    // Parse response
    var qlooResp TasteProfileResponse
    if err := json.Unmarshal(body, &qlooResp); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }

    // Check for API-level errors
    if qlooResp.Status != "success" && qlooResp.Status != "" {
        return nil, fmt.Errorf("Qloo API error: %s", qlooResp.Message)
    }

    return qlooResp.Segments, nil
}