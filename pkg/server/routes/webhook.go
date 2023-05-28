package routes

import (
	"net/http"

	"github.com/rookout/piper/pkg/conf"

	"github.com/gin-gonic/gin"
)

func addWebhookRoutes(cfg *conf.Config, rg *gin.RouterGroup) {
	health := rg.Group("/webhook")

	health.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, "healthy")
	})
}
