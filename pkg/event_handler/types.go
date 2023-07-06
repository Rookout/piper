package event_handler

import (
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"golang.org/x/net/context"
)

type EventNotifier interface {
	notify(ctx *context.Context, workflow *v1alpha1.Workflow) error
}
