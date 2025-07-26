package middleware

import (
    "fmt"
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/jesee-kuya/blue/internal/cache"
)

const (
    RateLimitMax     = 100
    RateLimitWindow  = time.Minute
    RateLimitPrefix  = "ratelimit:"
)

// RateLimitMiddleware creates a rate limiting middleware using Redis
func RateLimitMiddleware(redisClient *cache.RedisClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        clientIP := c.ClientIP()
        key := fmt.Sprintf("%s%s", RateLimitPrefix, clientIP)
        
        // Get current timestamp for rate limit window
        now := time.Now()
        windowStart := now.Truncate(RateLimitWindow).Unix()
        rateLimitKey := fmt.Sprintf("%s:%d", key, windowStart)

        // Increment counter for this IP in current window
        count, err := redisClient.Incr(rateLimitKey)
        if err != nil {
            // If Redis is down, allow the request but log error
            c.Header("X-RateLimit-Limit", strconv.Itoa(RateLimitMax))
            c.Header("X-RateLimit-Remaining", "0")
            c.Header("X-RateLimit-Reset", strconv.FormatInt(windowStart+int64(RateLimitWindow.Seconds()), 10))
            c.Next()
            return
        }

        // Set expiration on first increment
        if count == 1 {
            redisClient.Expire(rateLimitKey, RateLimitWindow)
        }

        // Calculate remaining requests and reset time
        remaining := RateLimitMax - int(count)
        if remaining < 0 {
            remaining = 0
        }
        resetTime := windowStart + int64(RateLimitWindow.Seconds())

        // Set rate limit headers
        c.Header("X-RateLimit-Limit", strconv.Itoa(RateLimitMax))
        c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
        c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

        // Check if rate limit exceeded
        if count > RateLimitMax {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "Rate limit exceeded",
                "retry_after": int64(RateLimitWindow.Seconds()) - (now.Unix() - windowStart),
            })
            c.Abort()
            return
        }

        c.Next()
    }
}