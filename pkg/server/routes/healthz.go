package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddHealthRoutes(rg *gin.RouterGroup) {
	health := rg.Group("/healthz")

	health.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, "healthy")
	})
}
