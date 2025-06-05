package auth

import (
	"net/http"
	"strings"
	"pubserver/internal/logger"

	"github.com/gin-gonic/gin"
)

const TokenTypeKey = "pubserver_token_type"

var tokens []TokenInfo

func SetupTokens(tokenList []TokenInfo) {
	tokens = tokenList
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			logger.Log.Warnw("Auth faltando Bearer", "ip", c.ClientIP())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token ausente"})
			return
		}
		tokenValue := strings.TrimPrefix(auth, "Bearer ")
		for _, t := range tokens {
			if tokenValue == t.Value && t.Active {
				logger.Log.Infow("Auth OK", "type", t.Type, "ip", c.ClientIP())
				c.Set(TokenTypeKey, t.Type)
				c.Next()
				return
			}
		}
		logger.Log.Warnw("Auth token inválido", "token", tokenValue, "ip", c.ClientIP())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
	}
}