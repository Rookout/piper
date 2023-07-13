package server

import (
	"github.com/rookout/piper/pkg/webhook_creator"
	"golang.org/x/net/context"
	"log"
	"net/http"
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

func (s *GracefulShutdown) Shutdown(httpServer *http.Server, webhookCreator *webhook_creator.WebhookCreatorImpl) {
	// Listen for the interrupt signal.
	<-s.ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	s.stop()

	log.Println("shutting down gracefully...")
	// The context is used to inform the server it has 10 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	webhookCreator.Stop(&ctx)

	err := httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

}
