package admin

import (
	"github.com/gin-gonic/gin"
	"pubserver/internal/auth"
	"pubserver/internal/db"
)

var DB *db.DB

func RegisterAdminRoutes(r *gin.Engine, database *db.DB) {
	DB = database

	r.GET("/admin/login", auth.AdminLoginHandler)
	r.POST("/admin/login", auth.AdminLoginHandler)
	r.GET("/admin/logout", auth.AdminLogoutHandler)

	admin := r.Group("/admin")
	admin.Use(auth.AdminSessionMiddleware())

	admin.GET("/", DashboardHandler)
	admin.GET("/tokens", TokensListHandler)
	admin.GET("/tokens/create", TokensCreateFormHandler)
	admin.POST("/tokens/create", TokensCreatePostHandler)
	admin.GET("/tokens/:id/toggle", TokenToggleHandler)
	admin.GET("/tokens/:id/delete", TokenDeleteHandler)
}