package openai

import (
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

// GetFunctionDefinitions returns the OpenAI function definitions for our available functions
func GetFunctionDefinitions() []openai.FunctionDefinition {
	return []openai.FunctionDefinition{
		{
			Name:        "search_marketplace",
			Description: "Search for products across multiple marketplaces (Amazon, eBay, Jumia) with optional price filtering",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"query": {
						Type:        jsonschema.String,
						Description: "Search query for products",
					},
					"min_price": {
						Type:        jsonschema.Number,
						Description: "Minimum price filter (optional)",
					},
					"max_price": {
						Type:        jsonschema.Number,
						Description: "Maximum price filter (optional)",
					},
				},
				Required: []string{"query"},
			},
		},
		{
			Name:        "get_taste_profile",
			Description: "Analyze a product description using Qloo's Taste AI to identify target audience segments with affinity scores",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"description": {
						Type:        jsonschema.String,
						Description: "Product description to analyze for audience segments",
					},
				},
				Required: []string{"description"},
			},
		},
		{
			Name:        "generate_ad_copy",
			Description: "Generate marketing copy and advertisements for a product targeting specific audience segments",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"product_title": {
						Type:        jsonschema.String,
						Description: "Title or name of the product",
					},
					"segments": {
						Type: jsonschema.Array,
						Items: &jsonschema.Definition{
							Type: jsonschema.String,
						},
						Description: "Target audience segments for the ad copy",
					},
				},
				Required: []string{"product_title", "segments"},
			},
		},
	}
}
