package middleware

import (
	"strings"

	"gst-api/pkg/response"

	"github.com/gin-gonic/gin"
)

const apiKeyHeader = "X-API-Key"

// APIKeyAuth returns a Gin middleware that validates requests using a static API key.
// The external company must include the header:  X-API-Key: <key>
func APIKeyAuth(expectedKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := strings.TrimSpace(c.GetHeader(apiKeyHeader))
		if key == "" {
			response.Unauthorized(c, "missing "+apiKeyHeader+" header")
			c.Abort()
			return
		}

		if key != expectedKey {
			response.Unauthorized(c, "invalid API key")
			c.Abort()
			return
		}

		c.Next()
	}
}
