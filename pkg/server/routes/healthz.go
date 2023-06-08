package routes

import (
	"net/http"

	"github.com/rookout/piper/pkg/conf"

	"github.com/gin-gonic/gin"
)

func AddHealthRoutes(cfg *conf.Config, rg *gin.RouterGroup) {
	health := rg.Group("/healthz")

	health.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, "healthy")
	})
}
