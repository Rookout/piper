package listener

import (
	"errors"
	"fmt"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/utils/strings/slices"
	"time"

	"k8s.io/client-go/tools/cache"
)

type K8sResourceEventsSubscriber struct {
	resource     runtime.Object
	namespace    string
	pubSub       PubSub
	watchStarted bool
	stopCh       chan struct{}
	restClient   *rest.RESTClient
}

const (
	ResourceCreated = "resource_created"
	ResourceUpdated = "resource_updated"
	ResourceDeleted = "resource_deleted"
)

var supportedEvents = []string{ResourceCreated, ResourceUpdated, ResourceDeleted}

func NewK8sResourceEventsSubscriber(resource runtime.Object, namespace string, restClient *rest.RESTClient) *K8sResourceEventsSubscriber {
	return &K8sResourceEventsSubscriber{
		resource:     resource,
		namespace:    namespace,
		pubSub:       NewSimplePubSub(),
		watchStarted: false,
		stopCh:       make(chan struct{}),
		restClient:   restClient,
	}
}

func (a *K8sResourceEventsSubscriber) Subscribe(eventName string, callback func(eventData any)) error {
	if !slices.Contains(supportedEvents, eventName) {
		return errors.New(fmt.Sprintf("invalid event - %s, can be one of %s", eventName, supportedEvents))
	}

	err := a.pubSub.Subscribe(eventName, callback)
	if err != nil {
		return err
	}

	if a.watchStarted {
		return nil
	}

	return a.startWatching()
}

func (a *K8sResourceEventsSubscriber) Stop() {
	close(a.stopCh)
}

func (a *K8sResourceEventsSubscriber) startWatching() error {
	// Create a new watcher for the specified resource in the specified namespace
	var resourceKind = "workflow"
	watcher := cache.NewListWatchFromClient(a.restClient, resourceKind, a.namespace, fields.Everything())

	// Start watching for events
	_, controller := cache.NewInformer(
		watcher,
		a.resource,
		time.Second*1,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				_ = a.pubSub.Publish(ResourceCreated, obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				_ = a.pubSub.Publish(ResourceUpdated, oldObj)
			},
			DeleteFunc: func(obj interface{}) {
				_ = a.pubSub.Publish(ResourceDeleted, obj)
			},
		},
	)

	// Run the controller
	go controller.Run(a.stopCh)

	a.watchStarted = true

	return nil
}
