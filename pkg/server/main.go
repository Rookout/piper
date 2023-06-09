package server

import (
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"golang.org/x/net/context"
	"log"
)

func Start(ctx context.Context, stop context.CancelFunc, cfg *conf.GlobalConfig, clients *clients.Clients) {

	srv := NewServer(cfg, clients)
	gracefulShutdownHandler := NewGracefulShutdown(ctx, stop)
	httpServer := srv.ListenAndServe()

	err := clients.GitProvider.SetWebhook()
	if err != nil {
		panic(err)
	}

	gracefulShutdownHandler.Shutdown(httpServer, clients)

	log.Println("Server exiting")
}
