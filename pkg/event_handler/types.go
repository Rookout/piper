package event_handler

import "k8s.io/apimachinery/pkg/watch"

type EventHandler interface {
	handle(workflowChan <-chan watch.Event)
}
