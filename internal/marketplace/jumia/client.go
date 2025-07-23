package jumia

import (
    "fmt"
    "io"
    "net/http"
    "net/url"
    "regexp"
    "strconv"
    "strings"
    "time"

    "github.com/jesee-kuya/blue/internal/marketplace"
)

// Client represents a Jumia web scraping client
type Client struct {
    baseURL    string
    httpClient *http.Client
}

// NewClient creates a new Jumia client
func NewClient() *Client {
    return &Client{
        baseURL: "https://www.jumia.com.ng",
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

// Search searches for products on Jumia using web scraping
func (c *Client) Search(query string, minPrice, maxPrice float64) ([]marketplace.Product, error) {
    // Build search URL
    searchURL := fmt.Sprintf("%s/catalog/?q=%s", c.baseURL, url.QueryEscape(query))
    
    // Create request
    req, err := http.NewRequest("GET", searchURL, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    // Set headers to mimic browser
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

    // Make request
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to make request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Jumia returned status %d", resp.StatusCode)
    }

    // Read response
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }

    // Parse HTML and extract products
    products, err := c.parseProducts(string(body))
    if err != nil {
        return nil, fmt.Errorf("failed to parse products: %w", err)
    }

    // Filter by price range
    var filteredProducts []marketplace.Product
    for _, product := range products {
        if (minPrice == 0 || product.Price >= minPrice) && 
           (maxPrice == 0 || product.Price <= maxPrice) {
            filteredProducts = append(filteredProducts, product)
        }
    }

    return filteredProducts, nil
}

// parseProducts extracts product information from HTML
func (c *Client) parseProducts(html string) ([]marketplace.Product, error) {
    var products []marketplace.Product

    // Regex patterns for extracting product data
    titlePattern := regexp.MustCompile(`<h3[^>]*class="[^"]*name[^"]*"[^>]*>([^<]+)</h3>`)
    pricePattern := regexp.MustCompile(`<div[^>]*class="[^"]*prc[^"]*"[^>]*>.*?₦\s*([0-9,]+)`)
    linkPattern := regexp.MustCompile(`<a[^>]*href="([^"]*)"[^>]*class="[^"]*core[^"]*"`)

    // Extract titles
    titleMatches := titlePattern.FindAllStringSubmatch(html, -1)
    
    // Extract prices
    priceMatches := pricePattern.FindAllStringSubmatch(html, -1)
    
    // Extract links
    linkMatches := linkPattern.FindAllStringSubmatch(html, -1)

    // Combine extracted data
    maxItems := len(titleMatches)
    if len(priceMatches) < maxItems {
        maxItems = len(priceMatches)
    }
    if len(linkMatches) < maxItems {
        maxItems = len(linkMatches)
    }

    for i := 0; i < maxItems; i++ {
        title := strings.TrimSpace(titleMatches[i][1])
        priceStr := strings.ReplaceAll(priceMatches[i][1], ",", "")
        price, err := strconv.ParseFloat(priceStr, 64)
        if err != nil {
            continue // Skip items with invalid prices
        }

        link := linkMatches[i][1]
        if !strings.HasPrefix(link, "http") {
            link = c.baseURL + link
        }

        products = append(products, marketplace.Product{
            Title: title,
            Price: price,
            Link:  link,
        })
    }

    return products, nil
}