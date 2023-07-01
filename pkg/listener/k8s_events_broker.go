package listener

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type K8sEventsBroker struct {
	resource  string
	namespace string
	broker    *EventBrokerExample
}

func NewK8sEventBroker(resource string, namespace string) *K8sEventsBroker {
	return &K8sEventsBroker{
		resource:  resource,
		namespace: namespace,
		broker:    NewEventBrokerExample(),
	}
}

func (a *K8sEventsBroker) Subscribe(eventName string, callback func(eventData any)) error {
	return a.broker.Subscribe(eventName, callback)
}

func (a *K8sEventsBroker) Publish(eventName string, eventData any) error {
	return a.broker.Publish(eventName, eventData)
}

func (a *K8sEventsBroker) Start() error {
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
				_ = a.Publish(fmt.Sprintf("%s_created", a.resource), obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				_ = a.Publish(fmt.Sprintf("%s_updated", a.resource), oldObj)
			},
			DeleteFunc: func(obj interface{}) {
				_ = a.Publish(fmt.Sprintf("%s_deleted", a.resource), obj)
			},
		},
	)

	// Run the controller
	stopCh := make(chan struct{})
	defer close(stopCh)
	go controller.Run(stopCh)

	// Wait indefinitely
	select {}
}
