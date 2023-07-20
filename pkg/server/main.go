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
	srv.Start()

	gracefulShutdownHandler.Shutdown(srv)

	log.Println("Server exiting")
}
