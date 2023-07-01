package listener

import (
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEventBrokerExample(t *testing.T) {

	var pubSub PubSub = NewEventPubSubExample()

	notifier := func(eventData any) {
		assert.Equal(t, "workflow completed", eventData.(v1alpha1.Workflow).Status.Message)
	}

	_ = pubSub.Subscribe("workflow_event", notifier)
	_ = pubSub.Publish("workflow_event", v1alpha1.Workflow{
		Status: v1alpha1.WorkflowStatus{
			Message: "workflow completed",
		},
	})
}
