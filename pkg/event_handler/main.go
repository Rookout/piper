package event_handler

import (
	"context"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
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
			{Key: "piper.rookout.com/notified",
				Operator: metav1.LabelSelectorOpNotIn,
				Values: []string{
					string(v1alpha1.WorkflowSucceeded),
					string(v1alpha1.WorkflowFailed),
					string(v1alpha1.WorkflowError),
				}}, // mean that there already completed and notified
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
	}()
	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
}
