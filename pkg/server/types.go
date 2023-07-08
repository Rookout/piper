package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/server/routes"
	"github.com/rs/zerolog"
)

type Server struct {
	router *gin.Engine
	logger *zerolog.Logger
	config *conf.GlobalConfig
}

func NewServer(config *conf.GlobalConfig, logger *zerolog.Logger) *Server {
	srv := &Server{
		router: gin.New(),
		logger: logger,
		config: config,
	}

	return srv
}

func (s *Server) registerMiddlewares() {
	s.router.Use(
		gin.LoggerWithConfig(gin.LoggerConfig{
			SkipPaths: []string{"/healthz"},
		}),
		gin.Recovery(),
	)

}

func (s *Server) getRoutes(cfg *conf.GlobalConfig, clients *clients.Clients, router *gin.Engine) {
	v1 := router.Group("/")
	routes.AddHealthRoutes(cfg, v1)
	routes.AddWebhookRoutes(cfg, clients, v1)
}
