package routes

import (
	"log"
	"net/http"

	webhook_hanlder "github.com/rookout/piper/pkg/webhook-hanlder"

	"github.com/rookout/piper/pkg/clients"

	"github.com/rookout/piper/pkg/conf"

	"github.com/gin-gonic/gin"
)

func AddWebhookRoutes(cfg *conf.Config, clients *clients.Clients, rg *gin.RouterGroup) {
	webhook := rg.Group("/webhook")

	webhook.POST("", func(c *gin.Context) {
		ctx := c.Copy().Request.Context()
		webhookPayload, err := clients.Git.HandlePayload(c.Request, []byte(cfg.GitConfig.WebhookSecret))
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if webhookPayload.Event == "ping" {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
			return
		}

		wh, err := webhook_hanlder.NewWebhookHandler(cfg, clients, webhookPayload)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Printf("failed to create webhook handler, error: %v", err)
			return
		}

		err = webhook_hanlder.HandleWebhook(&ctx, wh)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Printf("failed to handle webhook, error: %v", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
