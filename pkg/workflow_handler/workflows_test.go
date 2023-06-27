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
	returnConfigName, err := wfcImpl.SelectConfig(workflowsBatch)

	// Assert the expected output
	assert.Equal("default", returnConfigName)
	assert.Nil(err)

	// Test case 2
	configName = "config1"
	workflowsBatch = &common.WorkflowsBatch{
		Config:  &configName, // Set the desired config name here
		Payload: &git_provider.WebhookPayload{},
	}

	// Call the SelectConfig function
	returnConfigName, err = wfcImpl.SelectConfig(workflowsBatch)

	// Assert the expected output
	assert.Equal("config1", returnConfigName)
	assert.Nil(err)

	// Test case 3 - selection of non-existing config when default config exists
	configName = "notInConfigs"
	workflowsBatch = &common.WorkflowsBatch{
		Config:  &configName, // Set the desired config name here
		Payload: &git_provider.WebhookPayload{},
	}

	// Call the SelectConfig function
	returnConfigName, err = wfcImpl.SelectConfig(workflowsBatch)

	// Assert the expected output
	assert.Equal("default", returnConfigName)
	assert.NotNil(err)

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
	returnConfigName, err = wfcImpl4.SelectConfig(workflowsBatch)

	// Assert the expected output
	assert.NotNil(returnConfigName)
	assert.NotNil(err)
}

func TestCreateWorkflow(t *testing.T) {
	var wfc *conf.WorkflowsConfig
	var wfs *conf.WorkflowServerConfig

	// Create a WorkflowsClientImpl instance
	assert := assertion.New(t)
	// Create a mock WorkflowsClientImpl object with necessary dependencies
	wfc = &conf.WorkflowsConfig{Configs: map[string]*conf.ConfigInstance{
		"default": {Spec: v1alpha1.WorkflowSpec{},
			OnExit: []v1alpha1.DAGTask{}},
		"config1": {Spec: v1alpha1.WorkflowSpec{},
			OnExit: []v1alpha1.DAGTask{}},
	}}

	wfs = &conf.WorkflowServerConfig{Namespace: "default"}

	wfcImpl := &WorkflowsClientImpl{
		cfg: &conf.GlobalConfig{
			WorkflowsConfig:      *wfc,
			WorkflowServerConfig: *wfs,
		},
	}

	// Create a sample WorkflowSpec
	spec := &v1alpha1.WorkflowSpec{
		// Assign values to the fields of WorkflowSpec
		// ...

		// Example assignments:
		Entrypoint: "my-entrypoint",
	}

	// Create a sample WorkflowsBatch
	workflowsBatch := &common.WorkflowsBatch{
		Payload: &git_provider.WebhookPayload{
			Repo:   "my-repo",
			Branch: "my-branch",
			User:   "my-user",
			Commit: "my-commit",
		},
	}

	// Call the CreateWorkflow method
	workflow, err := wfcImpl.CreateWorkflow(spec, workflowsBatch)

	// Assert that no error occurred
	assert.NoError(err)

	// Assert that the returned workflow is not nil
	assert.NotNil(workflow)

	// Assert that the workflow's GenerateName, Namespace, and Labels are assigned correctly
	assert.Equal("my-repo-my-branch-", workflow.ObjectMeta.GenerateName)
	assert.Equal(wfcImpl.cfg.Namespace, workflow.ObjectMeta.Namespace)
	assert.Equal(map[string]string{
		"repo":   "my-repo",
		"branch": "my-branch",
		"user":   "my-user",
		"commit": "my-commit",
	}, workflow.ObjectMeta.Labels)

	// Assert that the workflow's Spec is assigned correctly
	assert.Equal(*spec, workflow.Spec)
}
