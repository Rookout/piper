package main

import (
	"log"

	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/git"
	"github.com/rookout/piper/pkg/server"
	"github.com/rookout/piper/pkg/utils"
	workflowHandler "github.com/rookout/piper/pkg/workflow-handler"

	rookout "github.com/Rookout/GoSDK"
)

func main() {
	cfg, err := conf.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load the configuration for Piper, error: %v", err)
	}

	if cfg.RookoutConfig.Token != "" {
		labels := utils.StringToMap(cfg.RookoutConfig.Labels)
		rookout.Start(rookout.RookOptions{Token: cfg.RookoutConfig.Token, Labels: labels})
	}

	git, err := git.NewGitProviderClient(cfg)
	if err != nil {
		log.Fatalf("failed to load the Git client for Piper, error: %v", err)
	}
	workflows, err := workflowHandler.NewWorkflowsClient(cfg)
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

	//err = common.Git.UnsetWebhook()
	//if err != nil {
	//	panic(err)
	//}

	server.Start(cfg, clients)
}
