package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jesee-kuya/blue/handler"
)

func main() {
	r := gin.Default()

	r.GET("/health", handler.HealthCheck)
	r.GET("/search", handler.SearchHandler)
	r.GET("/marketing", handler.MarketingHandler)

	r.Run(":8080")
}
