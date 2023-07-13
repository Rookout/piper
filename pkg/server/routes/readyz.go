package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddReadyRoutes(rg *gin.RouterGroup) {
	health := rg.Group("/readyz")

	health.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ready")
	})
}
