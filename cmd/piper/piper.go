package main

import (
	"log"

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
		log.Fatalf("failed to load the configuration for the piper, error: %v", err)
	}
}

func main() {
	server.Start(cfg)
}
