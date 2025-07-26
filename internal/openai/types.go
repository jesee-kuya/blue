package openai

// FunctionCall represents a function call request from OpenAI
type FunctionCall struct {
    Name      string                 `json:"name"`
    Arguments map[string]interface{} `json:"arguments"`
}

// FunctionResult represents the result of executing a function
type FunctionResult struct {
    FunctionName string      `json:"function_name"`
    Result       interface{} `json:"result"`
    Error        string      `json:"error,omitempty"`
}

// SearchMarketplaceArgs represents arguments for marketplace search function
type SearchMarketplaceArgs struct {
    Query    string  `json:"query"`
    MinPrice float64 `json:"min_price"`
    MaxPrice float64 `json:"max_price"`
}

// GetTasteProfileArgs represents arguments for taste profile function
type GetTasteProfileArgs struct {
    Description string `json:"description"`
}

// GenerateAdCopyArgs represents arguments for ad copy generation function
type GenerateAdCopyArgs struct {
    ProductTitle string   `json:"product_title"`
    Segments     []string `json:"segments"`
}

// AdCopyResult represents the result of ad copy generation
type AdCopyResult struct {
    Headlines    []string `json:"headlines"`
    Descriptions []string `json:"descriptions"`
    CallToAction string   `json:"call_to_action"`
}