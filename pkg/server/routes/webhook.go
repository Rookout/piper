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
			log.Fatalf("failed to create webhook handler, error: %v", err)
		}

		err = wh.RegisterTriggers()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Fatalf("failed to register triggers, error: %v", err)
			return
		} else {
			log.Printf("successfully registered triggers for repo: %s branch: %s", wh.Payload.Repo, wh.Payload.Branch)
		}

		err = wh.ExecuteMatchingTriggers()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Fatalf("failed to execute matching triggers, error: %v", err)
			return
		} else {
			log.Printf("successfully executed matching triggers for repo: %s branch: %s", wh.Payload.Repo, wh.Payload.Branch)
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
