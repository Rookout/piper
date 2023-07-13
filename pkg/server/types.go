package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/webhook_creator"
	"net/http"
)

type Server struct {
	router         *gin.Engine
	config         *conf.GlobalConfig
	clients        *clients.Clients
	webhookCreator *webhook_creator.WebhookCreatorImpl
}

type Interface interface {
	startServer() *http.Server
	registerMiddlewares()
	getRoutes()
	ListenAndServe() *http.Server
}
