package auth

import (
	"net/http"
	"pubserver/internal/db"
	"strings"

	"github.com/gin-gonic/gin"
)

const TokenTypeKey = "pubserver_token_type"

func AuthMiddlewareDB(database *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}
		tokenValue := strings.TrimPrefix(auth, "Bearer ")
		token, err := database.GetTokenValue(tokenValue)
		if err != nil || !token.Active {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		c.Set(TokenTypeKey, token.Type)
		c.Next()
	}
}