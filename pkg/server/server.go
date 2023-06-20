package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/server/routes"
	"log"
)

var router = gin.Default()

func Start(cfg *conf.GlobalConfig, clients *clients.Clients) {
	getRoutes(cfg, clients)
	err := router.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func getRoutes(cfg *conf.GlobalConfig, clients *clients.Clients) {
	v1 := router.Group("/")
	routes.AddHealthRoutes(cfg, v1)
	routes.AddWebhookRoutes(cfg, clients, v1)
}
