package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	webhookHandler "github.com/rookout/piper/pkg/webhook_handler"
)

func AddWebhookRoutes(cfg *conf.GlobalConfig, clients *clients.Clients, rg *gin.RouterGroup) {
	webhook := rg.Group("/webhook")

	webhook.POST("ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "pong"})
		return
	})

	webhook.POST("", func(c *gin.Context) {
		ctx := c.Copy().Request.Context()
		webhookPayload, err := clients.GitProvider.HandlePayload(c.Request, []byte(cfg.GitProviderConfig.WebhookSecret))
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if webhookPayload.Event == "ping" {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
			return
		}

		wh, err := webhookHandler.NewWebhookHandler(cfg, clients, webhookPayload)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Printf("failed to create webhook handler, error: %v", err)
			return
		}

		workflowsBatches, err := webhookHandler.HandleWebhook(&ctx, wh)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Printf("failed to handle webhook, error: %v", err)
			return
		}

		for _, wf := range workflowsBatches {
			err = clients.Workflows.HandleWorkflowBatch(&ctx, wf)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				log.Printf("failed to handle workflow, error: %v", err)
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
