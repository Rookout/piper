package routes

import (
	"github.com/rookout/piper/pkg/webhook_creator"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddHealthRoutes(rg *gin.RouterGroup, wc *webhook_creator.WebhookCreatorImpl) {
	health := rg.Group("/healthz")

	health.GET("", func(c *gin.Context) {
		ctx := c.Copy().Request.Context()
		err := wc.RunDiagnosis(&ctx)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, "healthy")
	})
}
