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

	notifier := NewGithubEventNotifier(cfg, clients)
	handler := &workflowEventHandler{
		Clients:  clients,
		Notifier: notifier,
	}
	go func() {
		for event := range watcher.ResultChan() {
			err = handler.Handle(ctx, &event)
			if err != nil {
				log.Printf("[event handler] failed to Handle workflow event %s", err) // ERROR
			}
		}
	}()
}
