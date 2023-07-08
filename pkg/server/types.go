package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
)

type Server struct {
	router  *gin.Engine
	config  *conf.GlobalConfig
	clients *clients.Clients
}
