package main

import (
	"html/template"
	"os"
	"pubserver/internal/admin"
	"pubserver/internal/api"
	"pubserver/internal/auth"
	"pubserver/internal/db"
	"pubserver/internal/gitlab"
	"pubserver/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func formatUnixTime(ts int64) string {
	return time.Unix(ts, 0).UTC().Format("2006-01-02 15:04")
}

func main() {
	logger.InitLogger()

	gitlabToken := os.Getenv("GITLAB_TOKEN")
	gitlabProject := os.Getenv("GITLAB_PROJECT")
	sqlitePath := os.Getenv("SQLITE_PATH")
	adminPw := os.Getenv("ADMIN_PASSWORD")
	if gitlabToken == "" || gitlabProject == "" || sqlitePath == "" || adminPw == "" {
		logger.Log.Fatalw("Set env variables: GITLAB_TOKEN, GITLAB_PROJECT, SQLITE_PATH, ADMIN_PASSWORD.")
	}

	dbase, err := db.Open(sqlitePath)
	if err != nil {
		logger.Log.Fatalw("Failed to open database", "err", err)
	}

	glClient, err := gitlab.NewClient(gitlabToken, gitlabProject)
	if err != nil {
		logger.Log.Fatalw("Failed to create GitLab client", "err", err)
	}

	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"formatUnixTime": formatUnixTime,
	})
	r.LoadHTMLGlob("web/templates/*.html")
	admin.RegisterAdminRoutes(r, dbase)

	// API protected by tokens from db
	r.Use(auth.AuthMiddlewareDB(dbase))
	api.RegisterPubRoutesWithStats(r, glClient, dbase)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Log.Infow("Server running", "port", port)
	r.Run(":" + port)
}