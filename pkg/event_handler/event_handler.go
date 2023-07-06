package event_handler

import (
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"k8s.io/apimachinery/pkg/watch"
	"log"
	"strconv"
)

type eventHandlerImpl struct{}

func (eh *eventHandlerImpl) handler(workflowChan <-chan watch.Event) {
	log.Printf("event handler started")
	for event := range workflowChan {
		wf, ok := event.Object.(*v1alpha1.Workflow)
		if !ok {
			log.Printf("Event object is not a Workflow object, it's kind is: %s", event.DeepCopy().Object.GetObjectKind())
			return
		}
		fmt.Printf(
			"evnet are: %s, %s phase: %s completed: %s, message: %s\n",
			event.Type,
			wf.GetName(),
			wf.Status.Phase,
			strconv.FormatBool(wf.Status.Phase.Completed()),
			wf.Status.Message,
		)
	}
}
