package event_handler

import (
	"context"
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"k8s.io/apimachinery/pkg/watch"
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
		clients:  clients,
		notifier: notifier,
	}
	go func() {
		for event := range watcher.ResultChan() {
			err = handler.handle(ctx, &event)
			if err != nil {
				log.Printf("[event handler] failed to handle workflow event %s", err) // ERROR
			}
		}
	}()
}

type workflowEventHandler struct {
	clients  *clients.Clients
	notifier EventNotifier
}

func (weh *workflowEventHandler) handle(ctx context.Context, event *watch.Event) error {
	workflow, ok := event.Object.(*v1alpha1.Workflow)
	if !ok {
		return fmt.Errorf("event object is not a Workflow object, it's kind is: %s\n", event.DeepCopy().Object.GetObjectKind())
	}

	currentPiperNotifyLabelStatus, ok := workflow.GetLabels()["piper/notify"]
	if !ok {
		return fmt.Errorf("workflow %s missing piper/notify label\n", workflow.GetName())

	}

	if currentPiperNotifyLabelStatus == string(workflow.Status.Phase) {
		log.Printf("workflow %s already informed for %s status. skiping... \n", workflow.GetName(), workflow.Status.Phase) //INFO
		return nil
	}

	ctx = context.Background()
	err := weh.notifier.notify(&ctx, workflow)
	if err != nil {
		return fmt.Errorf("failed to notify workflow to github, error:%s\n", err)

	}

	err = weh.clients.Workflows.UpdatePiperNotifyStatus(&ctx, workflow.GetName(), string(workflow.Status.Phase))
	if err != nil {
		return fmt.Errorf("error in workflow %s status patch: %s", workflow.GetName(), err)
	}
	log.Printf(
		"[event handler] done with event of type: %s for worklfow: %s phase: %s message: %s\n",
		event.Type,
		workflow.GetName(),
		workflow.Status.Phase,
		workflow.Status.Message) //INFO

	return nil
}
