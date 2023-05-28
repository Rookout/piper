package server

import (
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/server/routes"
)

func Start(cfg *conf.Config) {
	routes.Run(&cfg)
}
