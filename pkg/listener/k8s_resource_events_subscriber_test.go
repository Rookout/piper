package listener

import (
	"encoding/json"
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func restClient(testServer *httptest.Server) (*rest.RESTClient, error) {
	c, err := rest.RESTClientFor(&rest.Config{
		Host: testServer.URL,
		ContentConfig: rest.ContentConfig{
			GroupVersion:         &v1.SchemeGroupVersion,
			NegotiatedSerializer: scheme.Codecs.WithoutConversion(),
		},
		Username: "user",
		Password: "pass",
	})
	return c, err
}

type fakeHandler struct {
	serveHttp func(http.ResponseWriter, *http.Request)
}

func (f fakeHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	f.serveHttp(resp, req)
}

func TestK8sResourceEventsSubscriber(t *testing.T) {
	workflows := v1alpha1.WorkflowList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "workflow",
			APIVersion: "v1alpha1",
		},
		Items: v1alpha1.Workflows{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "workflow",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "workflow",
				},
				Spec: v1alpha1.WorkflowSpec{},
				Status: v1alpha1.WorkflowStatus{
					Phase: "pending",
				},
			},
		},
	}

	expectedBody, err := json.Marshal(workflows)
	assert.Nil(t, err)
	f := fakeHandler{
		serveHttp: func(response http.ResponseWriter, req *http.Request) {
			response.Header().Set("Content-Type", "application/json")
			response.WriteHeader(200)
			_, _ = response.Write(expectedBody)
		},
	}
	testServer := httptest.NewServer(&f)
	defer testServer.Close()
	client, err := restClient(testServer)

	var subscriber Subscriber = NewK8sResourceEventsSubscriber(&v1alpha1.Workflow{}, "default", client)

	err = subscriber.Subscribe(ResourceUpdated, func(event any) {
		fmt.Printf("workflow status: %s", event)
	})

	time.Sleep(time.Second * 10)

	assert.NotNil(t, err)
}
