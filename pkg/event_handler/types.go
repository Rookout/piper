package event_handler

import (
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/watch"
)

type EventHandler interface {
	Handle(ctx context.Context, event *watch.Event) error
}

type EventNotifier interface {
	Notify(ctx *context.Context, workflow *v1alpha1.Workflow) error
}
