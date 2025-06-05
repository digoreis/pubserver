package admin

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

func DashboardHandler(c *gin.Context) {
	recent, _ := DB.RecentStats(20)
	downloads, _ := DB.CountByEvent("download")
	publishes, _ := DB.CountByEvent("publish")
	type row struct {
		Event   string
		Package string
		Version string
		Time    string
	}
	var rows []row
	for _, e := range recent {
		rows = append(rows, row{
			Event:   e.Event,
			Package: e.Package,
			Version: e.Version,
			Time:    time.Unix(e.OccurredAt, 0).UTC().Format("2006-01-02 15:04:05"),
		})
	}
	c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{
		"recent":    rows,
		"downloads": downloads,
		"publishes": publishes,
	})
}