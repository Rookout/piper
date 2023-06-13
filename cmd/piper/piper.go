package main

import (
	workflow_handler "github.com/rookout/piper/pkg/workflow-handler"
	"log"

	"github.com/rookout/piper/pkg/server"

	"github.com/rookout/piper/pkg/clients"

	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/git"
)

func main() {
	cfg, err := conf.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load the configuration for Piper, error: %v", err)
	}
	git, err := git.NewGitProviderClient(cfg)
	if err != nil {
		log.Fatalf("failed to load the Git client for Piper, error: %v", err)
	}
	workflows, err := workflow_handler.NewWorkflowsClient(cfg)
	if err != nil {
		log.Fatalf("failed to load the Argo Workflows client for Piper, error: %v", err)
	}
	clients := &clients.Clients{
		Git:       git,
		Workflows: workflows,
	}

	err = clients.Git.SetWebhook()
	if err != nil {
		panic(err)
	}

	//err = clients.Git.UnsetWebhook()
	//if err != nil {
	//	panic(err)
	//}

	server.Start(cfg, clients)
}
