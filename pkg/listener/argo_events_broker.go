package listener

import (
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type ArgoEventsBroker struct {
}

func (a *ArgoEventsBroker) Subscribe(_ string, _ func(eventData any)) error {
	panic("not implemented")
}

func (a *ArgoEventsBroker) Publish(_ string, _ any) error {
	panic("not implemented")
}

func (a *ArgoEventsBroker) Start() {
	// get pod's service account credentials
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new watcher for Pods in the specified namespace
	watcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", corev1.NamespaceDefault, nil)

	// Start watching for events
	_, controller := cache.NewInformer(
		watcher,
		&corev1.Pod{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pod := obj.(*corev1.Pod)
				_ = a.Publish("pod_created", pod)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				pod := newObj.(*corev1.Pod)
				_ = a.Publish("pod_updated", pod)
			},
			DeleteFunc: func(obj interface{}) {
				pod := obj.(*corev1.Pod)
				_ = a.Publish("pod_deleted", pod)
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
