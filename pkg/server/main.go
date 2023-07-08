package server

import (
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"golang.org/x/net/context"
	"log"
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
