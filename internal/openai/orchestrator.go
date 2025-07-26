package openai

import (
	"context"
	"regexp"
	"strconv"
	"strings"
)

// MessageIntent represents the classified intent of a user message
type MessageIntent struct {
	Type        IntentType `json:"type"`
	Product     string     `json:"product,omitempty"`
	MinPrice    float64    `json:"min_price,omitempty"`
	MaxPrice    float64    `json:"max_price,omitempty"`
	Description string     `json:"description,omitempty"`
}

// IntentType represents different types of user intents
type IntentType string

const (
	IntentSearch    IntentType = "search"
	IntentMarketing IntentType = "marketing"
	IntentCombined  IntentType = "combined"
	IntentUnknown   IntentType = "unknown"
)

// OrchestratorResponse represents the final orchestrated response
type OrchestratorResponse struct {
	Message       string                `json:"message"`
	SearchResults *SearchResultsSummary `json:"search_results,omitempty"`
	Marketing     *MarketingCopy        `json:"marketing,omitempty"`
	Errors        []string              `json:"errors,omitempty"`
}

// SearchResultsSummary represents summarized search results
type SearchResultsSummary struct {
	Products []ProductSummary `json:"products"`
	Count    int              `json:"count"`
	Query    string           `json:"query"`
}

// ProductSummary represents a simplified product for responses
type ProductSummary struct {
	Title string  `json:"title"`
	Price float64 `json:"price"`
	Link  string  `json:"link"`
}

// MarketingCopy represents marketing content
type MarketingCopy struct {
	Headlines    []string `json:"headlines"`
	Descriptions []string `json:"descriptions"`
	CallToAction string   `json:"call_to_action"`
	Segments     []string `json:"target_segments"`
}

// ProcessMessage orchestrates the handling of user messages
func (c *Client) ProcessMessage(ctx context.Context, message string) (*OrchestratorResponse, error) {
	intent := c.classifyIntent(message)

	switch intent.Type {
	case IntentSearch:
		return c.handleSearchIntent(ctx, intent)
	case IntentMarketing:
		return c.handleMarketingIntent(ctx, intent)
	case IntentCombined:
		return c.handleCombinedIntent(ctx, intent)
	default:
		return c.handleUnknownIntent(ctx, message)
	}
}

// classifyIntent analyzes the user message to determine intent
func (c *Client) classifyIntent(message string) MessageIntent {
	message = strings.ToLower(message)

	// Search patterns
	searchPatterns := []string{
		`(?i)\b(find|search|show|list|get)\b.*\b(product|item|listing)`,
		`(?i)\b(find|search|show)\s+me\b`,
		`(?i)\bunder\s+\$?\d+`,
		`(?i)\bless\s+than\s+\$?\d+`,
		`(?i)\bbetween\s+\$?\d+.*\$?\d+`,
	}

	// Marketing patterns
	marketingPatterns := []string{
		`(?i)\b(marketing|advertis|ad|campaign|copy|promo)\b`,
		`(?i)\b(create|generate|suggest|make).*\b(ad|marketing|copy)\b`,
		`(?i)\btarget\s+(audience|segment)`,
	}

	// Combined patterns
	combinedPatterns := []string{
		`(?i)\b(find|search).*\b(and|then).*\b(marketing|ad|copy)\b`,
		`(?i)\b(marketing|ad).*\b(for|about).*\b(find|search)\b`,
	}

	hasSearch := c.matchesAnyPattern(message, searchPatterns)
	hasMarketing := c.matchesAnyPattern(message, marketingPatterns)
	hasCombined := c.matchesAnyPattern(message, combinedPatterns)

	intent := MessageIntent{}

	if hasCombined || (hasSearch && hasMarketing) {
		intent.Type = IntentCombined
	} else if hasMarketing {
		intent.Type = IntentMarketing
	} else if hasSearch {
		intent.Type = IntentSearch
	} else {
		intent.Type = IntentUnknown
	}

	// Extract product/description
	intent.Product = c.extractProduct(message)
	intent.Description = c.extractDescription(message)

	// Extract price range
	intent.MinPrice, intent.MaxPrice = c.extractPriceRange(message)

	return intent
}

// matchesAnyPattern checks if message matches any of the given patterns
func (c *Client) matchesAnyPattern(message string, patterns []string) bool {
	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, message); matched {
			return true
		}
	}
	return false
}

// extractProduct extracts product name from message
func (c *Client) extractProduct(message string) string {
	// Remove common prefixes and suffixes
	cleanMessage := regexp.MustCompile(`(?i)\b(find|search|show|get|for|about|create|generate|marketing|ad|copy)\b`).ReplaceAllString(message, "")
	cleanMessage = regexp.MustCompile(`(?i)\bunder\s+\$?\d+.*`).ReplaceAllString(cleanMessage, "")
	cleanMessage = strings.TrimSpace(cleanMessage)

	// Extract meaningful product terms
	words := strings.Fields(cleanMessage)
	var productWords []string

	for _, word := range words {
		if len(word) > 2 && !c.isStopWord(word) {
			productWords = append(productWords, word)
		}
	}

	if len(productWords) > 0 {
		return strings.Join(productWords[:min(3, len(productWords))], " ")
	}

	return ""
}

// extractDescription extracts product description for marketing
func (c *Client) extractDescription(message string) string {
	// For marketing intents, use the entire cleaned message as description
	if strings.Contains(strings.ToLower(message), "marketing") ||
		strings.Contains(strings.ToLower(message), "ad") {
		return c.extractProduct(message)
	}
	return ""
}

// extractPriceRange extracts min and max prices from message
func (c *Client) extractPriceRange(message string) (float64, float64) {
	var minPrice, maxPrice float64

	// Pattern: "under $X" or "less than $X"
	underPattern := regexp.MustCompile(`(?i)\b(under|less\s+than)\s+\$?(\d+(?:\.\d{2})?)`)
	if matches := underPattern.FindStringSubmatch(message); len(matches) > 2 {
		if price, err := strconv.ParseFloat(matches[2], 64); err == nil {
			maxPrice = price
		}
	}

	// Pattern: "between $X and $Y"
	betweenPattern := regexp.MustCompile(`(?i)\bbetween\s+\$?(\d+(?:\.\d{2})?)\s+(?:and|to)\s+\$?(\d+(?:\.\d{2})?)`)
	if matches := betweenPattern.FindStringSubmatch(message); len(matches) > 2 {
		if min, err := strconv.ParseFloat(matches[1], 64); err == nil {
			minPrice = min
		}
		if max, err := strconv.ParseFloat(matches[2], 64); err == nil {
			maxPrice = max
		}
	}

	return minPrice, maxPrice
}

// isStopWord checks if a word should be filtered out
func (c *Client) isStopWord(word string) bool {
	stopWords := map[string]bool{
		"the": true, "and": true, "or": true, "but": true, "in": true,
		"on": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "me": true, "my": true, "i": true,
		"you": true, "it": true, "is": true, "are": true, "was": true,
		"were": true, "be": true, "been": true, "have": true, "has": true,
		"had": true, "do": true, "does": true, "did": true, "will": true,
		"would": true, "could": true, "should": true, "may": true, "might": true,
	}
	return stopWords[strings.ToLower(word)]
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
