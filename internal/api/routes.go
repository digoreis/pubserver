package api

import (
	"github.com/gin-gonic/gin"
	"pubserver/internal/gitlab"
)

func RegisterPubRoutes(r *gin.Engine, glClient *gitlab.Client) {
	r.POST("/api/packages/versions/new", NewVersionHandler)
	r.POST("/api/packages/versions/newUpload", NewUploadHandler)
	r.POST("/api/packages/versions/newUploadFinish", NewUploadFinishHandler(glClient))
	r.GET("/api/packages/:package/versions", ListPackageVersionsHandler(glClient))
	r.GET("/api/packages/:package/versions/:version", DownloadPackageHandler(glClient))
}