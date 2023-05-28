package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rookout/piper/pkg/conf"
)

var router = gin.Default()

// Run will start the server
func Run(cfg *conf.Config) {
	getRoutes(cfg)
	router.Run()
}

// getRoutes will create our routes of our entire application
// this way every group of routes can be defined in their own file
// so this one won't be so messy
func getRoutes(cfg *conf.Config) {
	v1 := router.Group("/")
	addHealthRoutes(cfg, v1)
	addWebhookRoutes(cfg, v1)
}
