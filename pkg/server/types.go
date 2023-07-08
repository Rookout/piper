package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rs/zerolog"
	"net/http"
)

type Server struct {
	router  *gin.Engine
	logger  *zerolog.Logger
	config  *conf.GlobalConfig
	handler http.Handler
}

func NewServer(config *conf.GlobalConfig, logger *zerolog.Logger) (*Server, error) {
	srv := &Server{
		router: gin.New(),
		logger: logger,
		config: config,
	}

	return srv, nil
}

func (s *Server) registerMiddlewares() {
	s.router.Use(
		gin.LoggerWithConfig(gin.LoggerConfig{
			SkipPaths: []string{"/healthz"},
		}),
		gin.Recovery(),
	)

}
