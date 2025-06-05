package auth

import (
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
)

const adminSessionKey = "admin_logged_in"

func AdminSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := c.Cookie(adminSessionKey)
		if err != nil || session != "true" {
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

func AdminLoginHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "admin_login.html", gin.H{"error": ""})
		return
	}
	password := c.PostForm("password")
	adminPw := os.Getenv("ADMIN_PASSWORD")
	if password == adminPw && adminPw != "" {
		c.SetCookie(adminSessionKey, "true", 3600, "/", "", false, true)
		c.Redirect(http.StatusFound, "/admin")
		return
	}
	c.HTML(http.StatusOK, "admin_login.html", gin.H{"error": "Invalid password"})
}

func AdminLogoutHandler(c *gin.Context) {
	c.SetCookie(adminSessionKey, "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/admin/login")
}