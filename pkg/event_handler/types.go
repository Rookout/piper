package event_handler

import (
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"k8s.io/apimachinery/pkg/watch"
)

type EventHandler interface {
	handle(workflowChan <-chan watch.Event)
}

type EventNotifier interface {
	notify(workflow v1alpha1.Workflow)
}
