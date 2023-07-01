package listener

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type K8sResourceEventsSubscriber struct {
	resource  string
	namespace string
	pubSub    PubSub
	started   bool
}

func NewK8sResourceEventsSubscriber(resource string, namespace string) *K8sResourceEventsSubscriber {
	return &K8sResourceEventsSubscriber{
		resource:  resource,
		namespace: namespace,
		pubSub:    NewEventPubSubExample(),
		started:   false,
	}
}

func (a *K8sResourceEventsSubscriber) Subscribe(eventName string, callback func(eventData any)) error {
	err := a.pubSub.Subscribe(eventName, callback)
	if err != nil {
		return err
	}

	if a.started {
		return nil
	}

	return a.start()
}

func (a *K8sResourceEventsSubscriber) start() error {
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

	// Create a new watcher for Pods in the specified namespace
	watcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), a.resource, a.namespace, nil)

	// Start watching for events
	_, controller := cache.NewInformer(
		watcher,
		&corev1.Pod{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				_ = a.pubSub.Publish(fmt.Sprintf("%s_created", a.resource), obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				_ = a.pubSub.Publish(fmt.Sprintf("%s_updated", a.resource), oldObj)
			},
			DeleteFunc: func(obj interface{}) {
				_ = a.pubSub.Publish(fmt.Sprintf("%s_deleted", a.resource), obj)
			},
		},
	)

	// Run the controller
	stopCh := make(chan struct{})
	defer close(stopCh)
	go controller.Run(stopCh)

	a.started = true
	// Wait indefinitely
	select {}
}
