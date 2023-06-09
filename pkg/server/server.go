package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/server/routes"
	"log"
	"net/http"
)

func NewServer(config *conf.GlobalConfig, clients *clients.Clients) *Server {
	srv := &Server{
		router:  gin.New(),
		config:  config,
		clients: clients,
	}

	return srv
}

func (s *Server) startServer() *http.Server {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: s.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

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

func (s *Server) getRoutes() {
	v1 := s.router.Group("/")
	routes.AddHealthRoutes(v1)
	routes.AddWebhookRoutes(s.config, s.clients, v1)
}

func (s *Server) ListenAndServe() *http.Server {

	s.registerMiddlewares()

	s.getRoutes()

	return s.startServer()
}
