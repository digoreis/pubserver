package api

import (
	"net/http"
	"os"
	"pubserver/internal/gitlab"

	"github.com/gin-gonic/gin"
)

var (
	// Inicialize conforme sua config
	gitlabToken   = os.Getenv("GITLAB_TOKEN")
	gitlabProject = "grupo/meu-pacote"
	glClient, _   = gitlab.NewClient(gitlabToken, gitlabProject)
)

// Handler para publicar pacote
func PublishPackage(c *gin.Context) {
	version := c.PostForm("version")
	description := c.PostForm("description")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo não enviado"})
		return
	}
	// Salva arquivo temporariamente
	tmpPath := "/tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, tmpPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao salvar"})
		return
	}
	// Publica no GitLab
	err = glClient.PublishPackage(version, description, tmpPath)
	os.Remove(tmpPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao publicar " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Publicado com sucesso"})
}

// Listar versões
func ListPackageVersions(c *gin.Context) {
	versions, err := glClient.ListVersions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao listar"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"versions": versions})
}

// Download do asset
func DownloadPackage(c *gin.Context) {
	version := c.Param("version")
	filename := c.Param("filename")
	url, err := glClient.GetPackageAsset(version, filename)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset não encontrado"})
		return
	}
	c.Redirect(http.StatusFound, url) // Redireciona para o asset no GitLab
}