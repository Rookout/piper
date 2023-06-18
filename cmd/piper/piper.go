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
		err = rookout.Start(rookout.RookOptions{Token: cfg.RookoutConfig.Token, Labels: labels})
		if err != nil {
			log.Printf("failed to start Rookout, error: %v\n", err)
		}
	}

	err = cfg.WorkflowConfig.WorkflowsSpecLoad("/piper-config/..data")
	if err != nil {
		log.Fatalf("Failed to load workflow spec configuration, error: %v", err)
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
