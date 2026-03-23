package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger logs each incoming request with method, path, status and latency.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Printf("[%d] %s %s  latency=%s",
			c.Writer.Status(),
			c.Request.Method,
			c.Request.URL.Path,
			time.Since(start),
		)
	}
}
