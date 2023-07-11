package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"net/http"
)

type Server struct {
	router  *gin.Engine
	config  *conf.GlobalConfig
	clients *clients.Clients
}

type Interface interface {
	startServer() *http.Server
	registerMiddlewares()
	getRoutes()
	ListenAndServe() *http.Server
}

type Health interface {
	Check(msg healthCheck) error
	Recover(msg healthCheck) error
	Fail(msg healthCheck) error
	handler() error
}
