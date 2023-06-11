package routes

import (
	"net/http"

	"github.com/rookout/piper/pkg/clients"

	"github.com/rookout/piper/pkg/conf"

	"github.com/gin-gonic/gin"
)

func AddWebhookRoutes(cfg *conf.Config, clients *clients.Clients, rg *gin.RouterGroup) {
	webhook := rg.Group("/webhook")

	webhook.POST("", func(c *gin.Context) {
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.BindJSON(&json) == nil {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "json": json.Value})
		}
	})
}
