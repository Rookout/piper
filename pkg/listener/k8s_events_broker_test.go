package listener

import (
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestK8sEventsBroker(t *testing.T) {

	broker := NewK8sEventBroker("workflow", "default")

	_ = broker.Subscribe("workflow_updated", func(event any) {
		fmt.Printf("workflow status: %s", event.(v1alpha1.Workflow).Status.Message)
	})

	err := broker.Start()
	assert.NotNil(t, err)
}
