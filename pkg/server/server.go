package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/server/routes"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func Init() *gin.Engine {
	engine := gin.New()
	engine.Use(
		gin.LoggerWithConfig(gin.LoggerConfig{
			SkipPaths: []string{"/healthz"},
		}),
		gin.Recovery(),
	)
	return engine
}

func Start(cfg *conf.GlobalConfig, clients *clients.Clients) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := Init()

	getRoutes(cfg, clients, router)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

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

	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func getRoutes(cfg *conf.GlobalConfig, clients *clients.Clients, router *gin.Engine) {
	v1 := router.Group("/")
	routes.AddHealthRoutes(cfg, v1)
	routes.AddWebhookRoutes(cfg, clients, v1)
}

func gracefulShutdownHandler(ctx *context.Context, clients *clients.Clients) error {
	err := clients.GitProvider.UnsetWebhook(ctx)
	if err != nil {
		log.Println("Unset webhook error: ", err) // ERROR
		return err
	}

	return nil
}
