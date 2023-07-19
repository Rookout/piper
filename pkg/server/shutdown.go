package server

import (
	"golang.org/x/net/context"
	"log"
	"time"
)

type GracefulShutdown struct {
	ctx  context.Context
	stop context.CancelFunc
}

func NewGracefulShutdown(ctx context.Context, stop context.CancelFunc) *GracefulShutdown {
	return &GracefulShutdown{
		ctx:  ctx,
		stop: stop,
	}
}

func (s *GracefulShutdown) StopServices(ctx *context.Context, server *Server) {
	server.webhookCreator.Stop(ctx)
}

func (s *GracefulShutdown) Shutdown(server *Server) {
	// Listen for the interrupt signal.
	<-s.ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	s.stop()

	log.Println("shutting down gracefully...")
	// The context is used to inform the server it has 10 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.StopServices(&ctx, server)

	err := server.httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

}
