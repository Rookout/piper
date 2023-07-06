package event_handler

import (
	"context"
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
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
	go func() {
		for event := range watcher.ResultChan() {
			workflow, ok := event.Object.(*v1alpha1.Workflow)
			if !ok {
				log.Printf("event object is not a Workflow object, it's kind is: %s\n", event.DeepCopy().Object.GetObjectKind()) //ERROR
				return
			}
			ctx = context.Background()
			err = notifier.notify(&ctx, workflow)
			if err != nil {
				log.Printf("failed to notify workflow to github, error:%s\n", err)
				return
			}

			fmt.Printf(
				"evnet are: %s, %s phase: %s message: %s\n",
				event.Type,
				workflow.GetName(),
				workflow.Status.Phase,
				workflow.Status.Message,
			)
		}
	}()
}
