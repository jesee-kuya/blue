package ebay

// eBayResponse represents the response from eBay API
type eBayResponse struct {
    ItemSummaries []ItemSummary `json:"itemSummaries"`
    Total         int           `json:"total"`
}

// ItemSummary represents an item summary from eBay API
type ItemSummary struct {
    ItemID     string `json:"itemId"`
    Title      string `json:"title"`
    Price      Price  `json:"price"`
    ItemWebURL string `json:"itemWebUrl"`
}

// Price represents price information from eBay API
type Price struct {
    Value    string `json:"value"`
    Currency string `json:"currency"`
}