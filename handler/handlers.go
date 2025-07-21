package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(gin.HandlerFunc) {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}

func SearchHandler(gin.HandlerFunc ){
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func MarketingHandler(gin.HandlerFunc) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}