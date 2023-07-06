package event_handler

import "k8s.io/apimachinery/pkg/watch"

type EventHandler interface {
	handler(workflowChan <-chan watch.Event)
}
