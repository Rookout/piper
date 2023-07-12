package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/webhook_reconcile"
	"net/http"
)

type Server struct {
	router           *gin.Engine
	config           *conf.GlobalConfig
	clients          *clients.Clients
	webhookReconcile *webhook_reconcile.WebhookReconcileImpl
}

type Interface interface {
	startServer() *http.Server
	registerMiddlewares()
	getRoutes()
	ListenAndServe() *http.Server
}
