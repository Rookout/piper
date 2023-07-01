package listener

import (
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEventBrokerExample(t *testing.T) {

	broker := NewEventBrokerExample()

	notifier := func(eventData any) {
		assert.Equal(t, "workflow completed", eventData.(v1alpha1.Workflow).Status.Message)
	}

	_ = broker.Subscribe("workflow_event", notifier)
	_ = broker.Publish("workflow_event", v1alpha1.Workflow{
		Status: v1alpha1.WorkflowStatus{
			Message: "workflow completed",
		},
	})
}
