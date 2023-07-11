package routes

import (
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddHealthRoutes(cfg *conf.GlobalConfig, clients *clients.Clients, rg *gin.RouterGroup) {
	health := rg.Group("/healthz")

	health.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, "healthy")
	})
}
