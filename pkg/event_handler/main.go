package event_handler

import (
	"context"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"log"
)

func Start(cfg *conf.GlobalConfig, clients *clients.Clients) {
	ctx := context.Background()
	watcher, err := clients.Workflows.Watch(&ctx)
	if err != nil {
		log.Panicf("Failed to watch workflow error:%s", err)
	}

	handlerImpl := &eventHandlerImpl{}
	go func() {
		handlerImpl.handle(watcher.ResultChan())
	}()
}
