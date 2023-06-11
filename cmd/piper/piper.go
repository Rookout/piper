package main

import (
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

	clients := &clients.Clients{
		Git: git.NewGitProviderClient(cfg),
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
