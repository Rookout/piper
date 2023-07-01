package listener

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type K8sEventsBroker struct {
	resource  string
	namespace string
}

func NewK8sEventBroker(resource string, namespace string) *K8sEventsBroker {
	return &K8sEventsBroker{
		resource:  resource,
		namespace: namespace,
	}
}

func (a *K8sEventsBroker) Subscribe(_ string, _ func(eventData any)) error {
	panic("not implemented")
}

func (a *K8sEventsBroker) Publish(_ string, _ any) error {
	panic("not implemented")
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
				_ = a.Publish("pod_created", obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				_ = a.Publish("pod_updated", oldObj)
			},
			DeleteFunc: func(obj interface{}) {
				_ = a.Publish("pod_deleted", obj)
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
