func RegisterPubRoutesWithStats(r *gin.Engine, glClient *gitlab.Client, dbase *db.DB) {
    // ... outros endpoints ...

    r.GET("/api/packages/:package/versions/:version/interfaces", func(c *gin.Context) {
        pkg := c.Param("package")
        version := c.Param("version")

        // Busca a URL do asset interfaces.json na release do GitLab
        assetURL, err := glClient.GetReleaseAssetURL(pkg, version, "interfaces.json")
        if err != nil || assetURL == "" {
            c.JSON(http.StatusNotFound, gin.H{"error": "Interfaces not found"})
            return
        }

        // Faz o download do asset e retorna o JSON (pode também redirecionar)
        resp, err := http.Get(assetURL)
        if err != nil || resp.StatusCode != http.StatusOK {
            c.JSON(http.StatusNotFound, gin.H{"error": "Interfaces not found"})
            return
        }
        defer resp.Body.Close()
        c.Header("Content-Type", "application/json")
        c.Status(http.StatusOK)
        io.Copy(c.Writer, resp.Body)
    })
}