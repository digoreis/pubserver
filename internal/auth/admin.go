package auth

import (
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Admin-Token")
		adminToken := os.Getenv("ADMIN_TOKEN")
		if token != adminToken || adminToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Acesso admin negado"})
			return
		}
		c.Next()
	}
}