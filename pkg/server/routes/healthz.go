package routes

import (
	"github.com/rookout/piper/pkg/webhook_creator"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AddHealthRoutes(rg *gin.RouterGroup, wc *webhook_creator.WebhookCreatorImpl) {
	health := rg.Group("/healthz")

	health.GET("", func(c *gin.Context) {
		ctx := c.Copy().Request.Context()
		ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err := wc.RunDiagnosis(&ctx2)
		if err != nil {
			log.Printf("error from healthz endpint:%s\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, "healthy")
	})
}
