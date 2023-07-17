package event_handler

import (
	"context"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

func Start(ctx context.Context, stop context.CancelFunc, cfg *conf.GlobalConfig, clients *clients.Clients) {
	labelSelector := &metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{Key: "piper.rookout.com/notified",
				Operator: metav1.LabelSelectorOpExists},
		},
	}
	watcher, err := clients.Workflows.Watch(&ctx, labelSelector)
	if err != nil {
		log.Printf("[event handler] Failed to watch workflow error:%s", err)
		return
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
		log.Print("[event handler] stopped work, closing watcher")
		watcher.Stop()
		stop()
	}()
}
