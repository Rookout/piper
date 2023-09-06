package event_handler

import (
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rookout/piper/pkg/clients"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/watch"
	"log"
)

type workflowEventHandler struct {
	Clients  *clients.Clients
	Notifier EventNotifier
}

func (weh *workflowEventHandler) Handle(ctx context.Context, event *watch.Event) error {
	workflow, ok := event.Object.(*v1alpha1.Workflow)
	if !ok {
		return fmt.Errorf(
			"event object is not a Workflow object, got: %v\n",
			event.DeepCopy().Object,
		)
	}

	currentPiperNotifyLabelStatus, ok := workflow.GetLabels()["piper.rookout.com/notified"]
	if !ok {
		return fmt.Errorf(
			"workflow %s missing piper.rookout.com/notified label\n",
			workflow.GetName(),
		)
	}

	if currentPiperNotifyLabelStatus == string(workflow.Status.Phase) {
		log.Printf(
			"workflow %s already informed for %s status. skiping... \n",
			workflow.GetName(),
			workflow.Status.Phase,
		) //INFO
		return nil
	}

	err := weh.Notifier.Notify(&ctx, workflow)
	if err != nil {
		return fmt.Errorf("failed to Notify workflow to git provider, error:%s\n", err)
	}

	err = weh.Clients.Workflows.UpdatePiperWorkflowLabel(&ctx, workflow.GetName(), "notified", string(workflow.Status.Phase))
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
