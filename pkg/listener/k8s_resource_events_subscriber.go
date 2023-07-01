package listener

import (
	"errors"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/strings/slices"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type K8sResourceEventsSubscriber struct {
	resource     runtime.Object
	namespace    string
	pubSub       PubSub
	watchStarted bool
	stopCh       chan struct{}
}

const (
	ResourceCreated = "resource_created"
	ResourceUpdated = "resource_updated"
	ResourceDeleted = "resource_deleted"
)

var supportedEvents = []string{ResourceCreated, ResourceUpdated, ResourceDeleted}

func NewK8sResourceEventsSubscriber(resource runtime.Object, namespace string) *K8sResourceEventsSubscriber {
	return &K8sResourceEventsSubscriber{
		resource:     resource,
		namespace:    namespace,
		pubSub:       NewEventPubSubExample(),
		watchStarted: false,
		stopCh:       make(chan struct{}),
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
	// get pod's service account credentials
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	// Create a new Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// Create a new watcher for the specified resource in the specified namespace
	var resourceKind = a.resource.GetObjectKind().GroupVersionKind().Kind
	watcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), resourceKind, a.namespace, nil)

	// Start watching for events
	_, controller := cache.NewInformer(
		watcher,
		a.resource,
		time.Second*0,
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
