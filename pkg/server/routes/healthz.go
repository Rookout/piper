package routes

import (
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/webhook_reconcile"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddHealthRoutes(cfg *conf.GlobalConfig, clients *clients.Clients, rg *gin.RouterGroup, wr *webhook_reconcile.WebhookReconcileImpl) {
	health := rg.Group("/healthz")

	health.GET("", func(c *gin.Context) {
		//err := wr.RunTest()
		//if err != nil {
		//	c.JSON(http.StatusInternalServerError, err)
		//	return
		//}
		c.JSON(http.StatusOK, "healthy")
	})
}
