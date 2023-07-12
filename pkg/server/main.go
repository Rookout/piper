package server

import (
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"golang.org/x/net/context"
	"log"
)

func Start(ctx context.Context, stop context.CancelFunc, cfg *conf.GlobalConfig, clients *clients.Clients) {

	err := clients.GitProvider.SetWebhooks()
	if err != nil {
		panic(err)
	}
	srv := NewServer(cfg, clients)
	gracefulShutdownHandler := NewGracefulShutdown(ctx, stop)
	httpServer := srv.ListenAndServe()
	srv.webhookReconcile.ServeAndListen(ctx)

	gracefulShutdownHandler.Shutdown(httpServer, clients)

	log.Println("Server exiting")
}
