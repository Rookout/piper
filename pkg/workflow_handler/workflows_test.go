package workflow_handler

import (
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rookout/piper/pkg/common"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/git_provider"
	assertion "github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectConfig(t *testing.T) {
	var wfc *conf.WorkflowsConfig

	assert := assertion.New(t)
	// Create a sample WorkflowsBatch object for testing
	configName := "default"
	workflowsBatch := &common.WorkflowsBatch{
		Config:  &configName, // Set the desired config name here
		Payload: &git_provider.WebhookPayload{},
	}

	// Create a mock WorkflowsClientImpl object with necessary dependencies
	wfc = &conf.WorkflowsConfig{Configs: map[string]*conf.ConfigInstance{
		"default": {Spec: v1alpha1.WorkflowSpec{},
			OnExit: []v1alpha1.DAGTask{}},
		"config1": {Spec: v1alpha1.WorkflowSpec{},
			OnExit: []v1alpha1.DAGTask{}},
	}}

	wfcImpl := &WorkflowsClientImpl{
		cfg: &conf.GlobalConfig{
			WorkflowsConfig: *wfc,
		},
	}

	// Call the SelectConfig function
	returnConfigName := wfcImpl.SelectConfig(workflowsBatch)

	// Assert the expected output
	assert.Equal("default", returnConfigName)

	// Test case 2
	configName = "config1"
	workflowsBatch = &common.WorkflowsBatch{
		Config:  &configName, // Set the desired config name here
		Payload: &git_provider.WebhookPayload{},
	}

	// Call the SelectConfig function
	returnConfigName = wfcImpl.SelectConfig(workflowsBatch)

	// Assert the expected output
	assert.Equal("config1", returnConfigName)

	// Test case 3 - selection of non-existing config when default config exists
	configName = "notInConfigs"
	workflowsBatch = &common.WorkflowsBatch{
		Config:  &configName, // Set the desired config name here
		Payload: &git_provider.WebhookPayload{},
	}

	// Call the SelectConfig function
	returnConfigName = wfcImpl.SelectConfig(workflowsBatch)

	// Assert the expected output
	assert.Equal("default", returnConfigName)

	// Test case 4 - selection of non-existing config when default config not exists
	configName = "notInConfig"
	workflowsBatch = &common.WorkflowsBatch{
		Config:  &configName, // Set the desired config name here
		Payload: &git_provider.WebhookPayload{},
	}

	wfc4 := &conf.WorkflowsConfig{Configs: map[string]*conf.ConfigInstance{
		"config1": {Spec: v1alpha1.WorkflowSpec{},
			OnExit: []v1alpha1.DAGTask{}},
	}}

	wfcImpl4 := &WorkflowsClientImpl{
		cfg: &conf.GlobalConfig{
			WorkflowsConfig: *wfc4,
		},
	}

	// Call the SelectConfig function
	returnConfigName = wfcImpl4.SelectConfig(workflowsBatch)

	// Assert the expected output
	assert.NotNil(returnConfigName)

}
