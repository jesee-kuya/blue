package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jesee-kuya/blue/handler"
	"github.com/jesee-kuya/blue/internal/cache"
	"github.com/jesee-kuya/blue/internal/middleware"
)

func main() {
	r := gin.Default()

	// Initialize Redis client for rate limiting
	redisClient := cache.NewRedisClient()
	defer redisClient.Close()

	// Apply rate limiting middleware to protected routes
	rateLimited := r.Group("/")
	rateLimited.Use(middleware.RateLimitMiddleware(redisClient))
	{
		rateLimited.GET("/search", handler.SearchHandler)
		rateLimited.GET("/marketing", handler.MarketingHandler)
	}

	// Health check without rate limiting
	r.GET("/health", handler.HealthCheck)

	r.Run(":8080")
}
