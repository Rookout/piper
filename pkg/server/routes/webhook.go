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

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
