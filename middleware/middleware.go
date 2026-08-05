package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"maxwin/auth"
)

// AuthMiddleware verifies a per-user Bearer JWT and attaches claims to the context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, err := auth.BearerToken(c.GetHeader("Authorization"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing or invalid Authorization header",
			})
			c.Abort()
			return
		}

		claims, err := auth.ParseToken(raw)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set(auth.ContextUserIDKey, claims.UserID)
		c.Set(auth.ContextUsernameKey, claims.Username)
		c.Next()
	}
}

func UserID(c *gin.Context) string {
	v, _ := c.Get(auth.ContextUserIDKey)
	s, _ := v.(string)
	return s
}

func Username(c *gin.Context) string {
	v, _ := c.Get(auth.ContextUsernameKey)
	s, _ := v.(string)
	return s
}
