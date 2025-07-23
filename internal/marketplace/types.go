package marketplace

// Product represents a standardized product from any marketplace
type Product struct {
	Title string  `json:"title"`
	Price float64 `json:"price"`
	Link  string  `json:"link"`
}

// Client defines the interface that all marketplace clients must implement
type Client interface {
	Search(query string, minPrice, maxPrice float64) ([]Product, error)
}
