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
				log.Printf("[event handler] event object is not a Workflow object, it's kind is: %s\n", event.DeepCopy().Object.GetObjectKind()) //ERROR
				continue
			}

			currentPiperNotifyLabelStatus, ok := workflow.GetLabels()["piper/notify"]
			if !ok {
				log.Printf("[event handler] workflow %s missing piper/notify label\n", workflow.GetName()) //ERROR
				continue
			}

			if currentPiperNotifyLabelStatus == string(workflow.Status.Phase) {
				log.Printf("[event handler] workflow %s already informed for %s status. skiping... \n", workflow.GetName(), workflow.Status.Phase) //Info
				continue
			}

			ctx = context.Background()
			err = notifier.notify(&ctx, workflow)
			if err != nil {
				log.Printf("[event handler] failed to notify workflow to github, error:%s\n", err) //ERROR
				continue
			}

			err = clients.Workflows.UpdatePiperNotifyStatus(&ctx, workflow.GetName(), string(workflow.Status.Phase))
			if err != nil {
				log.Printf("[event handler] error in workflow %s status patch: %s", workflow.GetName(), err) //ERROR
			}
			fmt.Printf(
				"[event handler] evnet are: %s, %s phase: %s message: %s\n",
				event.Type,
				workflow.GetName(),
				workflow.Status.Phase,
				workflow.Status.Message) //INFO
		}
	}()
}
