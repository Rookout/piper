package server

import (
	"github.com/rookout/piper/pkg/clients"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"
)

type gracefulShutdown struct {
	ctx  context.Context
	stop context.CancelFunc
}

func NewGracefulShutdown(ctx context.Context, stop context.CancelFunc) *gracefulShutdown {
	return &gracefulShutdown{
		ctx:  ctx,
		stop: stop,
	}
}

func (s *gracefulShutdown) Shutdown(httpServer *http.Server, clients *clients.Clients) {
	// Listen for the interrupt signal.
	<-s.ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	s.stop()

	log.Println("shutting down gracefully...")
	// The context is used to inform the server it has 10 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := clients.GitProvider.UnsetWebhooks(&ctx)
	if err != nil {
		log.Println("Unset webhook error: ", err) // ERROR
	}

	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

}
