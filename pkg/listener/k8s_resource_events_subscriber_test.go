package listener

import (
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestK8sResourceEventsSubscriber(t *testing.T) {
	var subscriber Subscriber = NewK8sResourceEventsSubscriber("workflow", "default")

	err := subscriber.Subscribe("workflow_updated", func(event any) {
		fmt.Printf("workflow status: %s", event.(v1alpha1.Workflow).Status.Message)
	})

	assert.NotNil(t, err)
}
