package main

import (
	"log"

	"github.com/rookout/piper/pkg/clients"

	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/git"
	"github.com/rookout/piper/pkg/server"
)

func main() {
	cfg, err := conf.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load the configuration for Piper, error: %v", err)
	}

	clients := clients.Clients{
		Git: git.NewGitProviderClient(cfg),
	}

	server.Start(cfg, &clients)
}
