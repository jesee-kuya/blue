package openai

import (
    "fmt"
    "strings"

    "github.com/jesee-kuya/blue/internal/marketplace"
)

// convertSearchResults converts function call results to SearchResultsSummary
func (c *Client) convertSearchResults(result any, query string) *SearchResultsSummary {
    resultMap, ok := result.(map[string]any)
    if !ok {
        return &SearchResultsSummary{Query: query, Count: 0}
    }

    productsInterface, ok := resultMap["products"]
    if !ok {
        return &SearchResultsSummary{Query: query, Count: 0}
    }

    var products []ProductSummary
    
    // Handle different product formats
    switch p := productsInterface.(type) {
    case []marketplace.Product:
        for _, product := range p {
            products = append(products, ProductSummary{
                Title: product.Title,
                Price: product.Price,
                Link:  product.Link,
            })
        }
    case []any:
        for _, item := range p {
            if productMap, ok := item.(map[string]any); ok {
                product := ProductSummary{}
                if title, ok := productMap["title"].(string); ok {
                    product.Title = title
                }
                if price, ok := productMap["price"].(float64); ok {
                    product.Price = price
                }
                if link, ok := productMap["link"].(string); ok {
                    product.Link = link
                }
                products = append(products, product)
            }
        }
    }

    return &SearchResultsSummary{
        Products: products,
        Count:    len(products),
        Query:    query,
    }
}

// convertMarketingResults converts function call results to MarketingCopy
func (c *Client) convertMarketingResults(result any, segments []string) *MarketingCopy {
    switch r := result.(type) {
    case AdCopyResult:
        return &MarketingCopy{
            Headlines:    r.Headlines,
            Descriptions: r.Descriptions,
            CallToAction: r.CallToAction,
            Segments:     segments,
        }
    case map[string]any:
        marketing := &MarketingCopy{Segments: segments}
        
        if headlines, ok := r["headlines"].([]string); ok {
            marketing.Headlines = headlines
        } else if headlinesInterface, ok := r["headlines"].([]any); ok {
            for _, h := range headlinesInterface {
                if headline, ok := h.(string); ok {
                    marketing.Headlines = append(marketing.Headlines, headline)
                }
            }
        }
        
        if descriptions, ok := r["descriptions"].([]string); ok {
            marketing.Descriptions = descriptions
        } else if descriptionsInterface, ok := r["descriptions"].([]any); ok {
            for _, d := range descriptionsInterface {
                if description, ok := d.(string); ok {
                    marketing.Descriptions = append(marketing.Descriptions, description)
                }
            }
        }
        
        if cta, ok := r["call_to_action"].(string); ok {
            marketing.CallToAction = cta
        }
        
        return marketing
    }
    
    return &MarketingCopy{Segments: segments}
}

// extractSegments extracts segment names from taste profile results
func (c *Client) extractSegments(result any) []string {
    resultMap, ok := result.(map[string]any)
    if !ok {
        return []string{"General Consumers"}
    }

    segmentsInterface, ok := resultMap["segments"]
    if !ok {
        return []string{"General Consumers"}
    }

    var segments []string
    
    switch s := segmentsInterface.(type) {
    case []any:
        for _, item := range s {
            if segmentMap, ok := item.(map[string]any); ok {
                if name, ok := segmentMap["name"].(string); ok {
                    segments = append(segments, name)
                }
            } else if segmentStr, ok := item.(string); ok {
                segments = append(segments, segmentStr)
            }
        }
    case []string:
        segments = s
    }

    if len(segments) == 0 {
        return []string{"General Consumers"}
    }

    return segments
}

// formatSearchMessage creates a user-friendly search results message
func (c *Client) formatSearchMessage(results *SearchResultsSummary) string {
    if results.Count == 0 {
        return fmt.Sprintf("I couldn't find any products matching '%s'. Try adjusting your search terms or price range.", results.Query)
    }

    var message strings.Builder
    message.WriteString(fmt.Sprintf("I found %d products for '%s':\n\n", results.Count, results.Query))

    for i, product := range results.Products {
        if i >= 5 { // Limit to top 5 results in message
            message.WriteString(fmt.Sprintf("... and %d more results\n", results.Count-5))
            break
        }
        message.WriteString(fmt.Sprintf("• %s - $%.2f\n", product.Title, product.Price))
    }

    return message.String()
}

// formatMarketingMessage creates a user-friendly marketing copy message
func (c *Client) formatMarketingMessage(marketing *MarketingCopy, productName string) string {
    var message strings.Builder
    
    if productName != "" {
        message.WriteString(fmt.Sprintf("Here's marketing copy for '%s':\n\n", productName))
    } else {
        message.WriteString("Here's your marketing copy:\n\n")
    }

    if len(marketing.Segments) > 0 {
        message.WriteString(fmt.Sprintf("**Target Audience:** %s\n\n", strings.Join(marketing.Segments, ", ")))
    }

    if len(marketing.Headlines) > 0 {
        message.WriteString("**Headlines:**\n")
        for _, headline := range marketing.Headlines {
            message.WriteString(fmt.Sprintf("• %s\n", headline))
        }
        message.WriteString("\n")
    }

    if len(marketing.Descriptions) > 0 {
        message.WriteString("**Descriptions:**\n")
        for _, desc := range marketing.Descriptions {
            message.WriteString(fmt.Sprintf("• %s\n", desc))
        }
        message.WriteString("\n")
    }

    if marketing.CallToAction != "" {
        message.WriteString(fmt.Sprintf("**Call to Action:** %s\n", marketing.CallToAction))
    }

    return message.String()
}

// formatCombinedMessage creates a comprehensive message for combined results
func (c *Client) formatCombinedMessage(searchResults *SearchResultsSummary, marketing *MarketingCopy, productName string) string {
    var message strings.Builder

    if productName != "" {
        message.WriteString(fmt.Sprintf("Here's what I found for '%s':\n\n", productName))
    }

    // Add search results section
    if searchResults != nil && searchResults.Count > 0 {
        message.WriteString("## Product Listings\n")
        message.WriteString(fmt.Sprintf("Found %d products:\n", searchResults.Count))
        
        for i, product := range searchResults.Products {
            if i >= 3 { // Limit to top 3 for combined view
                message.WriteString(fmt.Sprintf("... and %d more\n", searchResults.Count-3))
                break
            }
            message.WriteString(fmt.Sprintf("• %s - $%.2f\n", product.Title, product.Price))
        }
        message.WriteString("\n")
    }

    // Add marketing section
    if marketing != nil && len(marketing.Headlines) > 0 {
        message.WriteString("## Marketing Copy\n")
        
        if len(marketing.Segments) > 0 {
            message.WriteString(fmt.Sprintf("**Target Audience:** %s\n\n", strings.Join(marketing.Segments, ", ")))
        }

        if len(marketing.Headlines) > 0 {
            message.WriteString("**Top Headlines:**\n")
            for i, headline := range marketing.Headlines {
                if i >= 2 { // Show top 2 headlines
                    break
                }
                message.WriteString(fmt.Sprintf("• %s\n", headline))
            }
            message.WriteString("\n")
        }

        if marketing.CallToAction != "" {
            message.WriteString(fmt.Sprintf("**Call to Action:** %s\n", marketing.CallToAction))
        }
    }

    if message.Len() == 0 {
        return "I encountered some issues processing your request. Please try again with more specific details."
    }

    return message.String()
}