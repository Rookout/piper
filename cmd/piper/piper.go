package main

import (
	"log"

	"github.com/rookout/piper/pkg/git"

	"github.com/rookout/piper/pkg/conf"
	server "github.com/rookout/piper/pkg/server"
)

var (
	cfg *conf.Config
)

func init() {
	var err error
	cfg, err = conf.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load the configuration for Piper, error: %v", err)
	}
	gitClient := git.NewGithubClient(cfg)
	gitClient.SetWebhook()
}

func main() {
	server.Start(cfg)
}