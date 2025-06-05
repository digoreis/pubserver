package admin

import (
	"net/http"
	"pubserver/internal/db"
	"strings"

	"github.com/gin-gonic/gin"
)

func TokensListHandler(c *gin.Context) {
	tokens, _ := DB.ListTokens()
	c.HTML(http.StatusOK, "admin_tokens.html", gin.H{
		"tokens": tokens,
		"new":    c.Query("new"),
	})
}

func TokensCreateFormHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_token_create.html", gin.H{"error": ""})
}

func TokensCreatePostHandler(c *gin.Context) {
	typ := strings.ToLower(c.PostForm("type"))
	desc := c.PostForm("description")
	if typ != "ci" && typ != "user" {
		c.HTML(http.StatusOK, "admin_token_create.html", gin.H{"error": "Invalid type"})
		return
	}
	token, err := DB.CreateToken(typ, desc)
	if err != nil {
		c.HTML(http.StatusOK, "admin_token_create.html", gin.H{"error": "Failed to create token"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/tokens?new="+token.Value)
}

func TokenToggleHandler(c *gin.Context) {
	id := c.Param("id")
	DB.ToggleToken(id)
	c.Redirect(http.StatusSeeOther, "/admin/tokens")
}

func TokenDeleteHandler(c *gin.Context) {
	id := c.Param("id")
	DB.DeleteToken(id)
	c.Redirect(http.StatusSeeOther, "/admin/tokens")
}