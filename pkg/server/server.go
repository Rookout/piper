package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/server/routes"
	"log"
)

func Init() *gin.Engine {
	engine := gin.New()
	engine.Use(
		gin.LoggerWithConfig(gin.LoggerConfig{
			SkipPaths: []string{"/healthz"},
		}),
		gin.Recovery(),
	)
	return engine
}

func Start(cfg *conf.GlobalConfig, clients *clients.Clients) {
	router := Init()
	getRoutes(cfg, clients, router)

	err := router.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func getRoutes(cfg *conf.GlobalConfig, clients *clients.Clients, router *gin.Engine) {
	v1 := router.Group("/")
	routes.AddHealthRoutes(cfg, v1)
	routes.AddWebhookRoutes(cfg, clients, v1)
}
