package main

import (
	rookout "github.com/Rookout/GoSDK"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/event_handler"
	"github.com/rookout/piper/pkg/git_provider"
	"github.com/rookout/piper/pkg/server"
	"github.com/rookout/piper/pkg/utils"
	workflowHandler "github.com/rookout/piper/pkg/workflow_handler"
	"golang.org/x/net/context"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := conf.LoadConfig()
	if err != nil {
		log.Panicf("failed to load the configuration for Piper, error: %v", err)
	}

	if cfg.RookoutConfig.Token != "" {
		labels := utils.StringToMap(cfg.RookoutConfig.Labels)
		err = rookout.Start(rookout.RookOptions{Token: cfg.RookoutConfig.Token, Labels: labels})
		if err != nil {
			log.Printf("failed to start Rookout, error: %v\n", err)
		}
	}

	err = cfg.WorkflowsConfig.WorkflowsSpecLoad("/piper-config/..data")
	if err != nil {
		log.Panicf("Failed to load workflow spec configuration, error: %v", err)
	}

	gitProvider, err := git_provider.NewGitProviderClient(cfg)
	if err != nil {
		log.Panicf("failed to load the Git client for Piper, error: %v", err)
	}
	workflows, err := workflowHandler.NewWorkflowsClient(cfg)
	if err != nil {
		log.Panicf("failed to load the Argo Workflows client for Piper, error: %v", err)
	}

	globalClients := &clients.Clients{
		GitProvider: gitProvider,
		Workflows:   workflows,
	}

	err = globalClients.GitProvider.SetWebhook()
	if err != nil {
		panic(err)
	}

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	event_handler.Start(ctx, stop, cfg, globalClients)
	server.Start(ctx, stop, cfg, globalClients)
}
