package listener

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestK8sEventsBroker(t *testing.T) {

	broker := NewK8sEventBroker("workflow", "default")

	err := broker.Start()
	assert.NotNil(t, err)

}
