package middleware

import (
	"strings"

	"gst-api/pkg/jwt"
	"gst-api/pkg/response"

	"github.com/gin-gonic/gin"
)

// JWTAuth validates the Bearer token from the Authorization header.
// Usage: Authorization: Bearer <token>
func JWTAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "missing Authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Unauthorized(c, "Authorization header format must be: Bearer <token>")
			c.Abort()
			return
		}

		claims, err := jwt.Validate(parts[1], jwtSecret)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token: "+err.Error())
			c.Abort()
			return
		}

		// Store claims in context so handlers can read them if needed
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
