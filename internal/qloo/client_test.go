package qloo

import (
    "net/http"
    "net/http/httptest"
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestQlooClient_GetTasteProfile_Success(t *testing.T) {
    // Mock Qloo API response
    mockResponse := `{
        "status": "success",
        "segments": [
            {
                "name": "Tech Enthusiasts",
                "affinity_score": 0.85
            },
            {
                "name": "Early Adopters",
                "affinity_score": 0.72
            },
            {
                "name": "Gaming Community",
                "affinity_score": 0.68
            }
        ]
    }`

    // Create test server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        assert.Equal(t, "POST", r.Method)
        assert.Equal(t, "/taste/profile", r.URL.Path)
        assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))
        assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(mockResponse))
    }))
    defer server.Close()

    // Create client with test server URL
    client := NewClientWithKey("test-api-key")
    client.baseURL = server.URL

    // Test GetTasteProfile
    segments, err := client.GetTasteProfile("High-performance gaming laptop with RGB lighting")

    assert.NoError(t, err)
    assert.Len(t, segments, 3)
    assert.Equal(t, "Tech Enthusiasts", segments[0].Name)
    assert.Equal(t, 0.85, segments[0].AffinityScore)
    assert.Equal(t, "Early Adopters", segments[1].Name)
    assert.Equal(t, 0.72, segments[1].AffinityScore)
}

func TestQlooClient_GetTasteProfile_EmptyDescription(t *testing.T) {
    client := NewClientWithKey("test-api-key")
    
    segments, err := client.GetTasteProfile("")
    
    assert.NoError(t, err)
    assert.Empty(t, segments)
}

func TestQlooClient_GetTasteProfile_NoAPIKey(t *testing.T) {
    client := NewClientWithKey("")
    
    segments, err := client.GetTasteProfile("test description")
    
    assert.Error(t, err)
    assert.Nil(t, segments)
    assert.Contains(t, err.Error(), "QLOO_API_KEY not set")
}

func TestQlooClient_GetTasteProfile_APIError(t *testing.T) {
    // Create test server that returns error
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte(`{"status": "error", "message": "Invalid API key"}`))
    }))
    defer server.Close()

    client := NewClientWithKey("invalid-key")
    client.baseURL = server.URL

    segments, err := client.GetTasteProfile("test description")

    assert.Error(t, err)
    assert.Nil(t, segments)
    assert.Contains(t, err.Error(), "Qloo API returned status 401")
}

func TestQlooClient_GetTasteProfile_InvalidJSON(t *testing.T) {
    // Create test server that returns invalid JSON
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`invalid json`))
    }))
    defer server.Close()

    client := NewClientWithKey("test-api-key")
    client.baseURL = server.URL

    segments, err := client.GetTasteProfile("test description")

    assert.Error(t, err)
    assert.Nil(t, segments)
    assert.Contains(t, err.Error(), "failed to parse response")
}

func TestQlooClient_GetTasteProfile_APILevelError(t *testing.T) {
    // Mock API response with error status
    mockResponse := `{
        "status": "error",
        "message": "Description too short",
        "segments": []
    }`

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(mockResponse))
    }))
    defer server.Close()

    client := NewClientWithKey("test-api-key")
    client.baseURL = server.URL

    segments, err := client.GetTasteProfile("short")

    assert.Error(t, err)
    assert.Nil(t, segments)
    assert.Contains(t, err.Error(), "Description too short")
}

func TestQlooClient_GetTasteProfile_NoSegments(t *testing.T) {
    // Mock response with no segments
    mockResponse := `{
        "status": "success",
        "segments": []
    }`

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(mockResponse))
    }))
    defer server.Close()

    client := NewClientWithKey("test-api-key")
    client.baseURL = server.URL

    segments, err := client.GetTasteProfile("generic product")

    assert.NoError(t, err)
    assert.Empty(t, segments)
}

func TestNewClient_WithEnvironmentVariable(t *testing.T) {
    // Set environment variable
    os.Setenv("QLOO_API_KEY", "env-api-key")
    defer os.Unsetenv("QLOO_API_KEY")

    client := NewClient()

    assert.Equal(t, "env-api-key", client.apiKey)
    assert.Equal(t, "https://api.qloo.com/v1", client.baseURL)
    assert.NotNil(t, client.httpClient)
}

func TestNewClient_WithoutEnvironmentVariable(t *testing.T) {
    // Ensure environment variable is not set
    os.Unsetenv("QLOO_API_KEY")

    client := NewClient()

    assert.Empty(t, client.apiKey)
    assert.Equal(t, "https://api.qloo.com/v1", client.baseURL)
    assert.NotNil(t, client.httpClient)
}