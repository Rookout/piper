package routes

import (
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddHealthRoutes(cfg *conf.GlobalConfig, clients *clients.Clients, rg *gin.RouterGroup) {
	health := rg.Group("/healthz")

	health.GET("", func(c *gin.Context) {
		ctx := c.Copy().Request.Context()
		err := clients.GitProvider.PingHooks(&ctx)
		if err != nil {
			log.Printf("error in pinging hooks %s", err)
			//err = clients.GitProvider.SetWebhook()
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, "healthy")
	})
}
