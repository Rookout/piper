package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/server/routes"
	"log"
	"net/http"
	"time"
)

func Start(ctx context.Context, stop context.CancelFunc, cfg *conf.GlobalConfig, clients *clients.Clients) {

	srv := NewServer(cfg, clients)
	httpServer := srv.ListenAndServe()

	err := clients.GitProvider.SetWebhook()
	if err != nil {
		panic(err)
	}

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully...")

	// The context is used to inform the server it has 10 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = gracefulShutdownHandler(&ctx, clients)

	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func gracefulShutdownHandler(ctx *context.Context, clients *clients.Clients) error {
	err := clients.GitProvider.UnsetWebhook(ctx)
	if err != nil {
		log.Println("Unset webhook error: ", err) // ERROR
		return err
	}

	return nil
}

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
