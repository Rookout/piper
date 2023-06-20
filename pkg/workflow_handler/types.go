package workflow_handler

import (
	"context"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rookout/piper/pkg/common"
)

type WorkflowsClient interface {
	ConstructTemplates(workflowsBatch *common.WorkflowsBatch, configName string) ([]v1alpha1.Template, error)
	ConstructSpec(templates []v1alpha1.Template, params []v1alpha1.Parameter, configName string) (*v1alpha1.WorkflowSpec, error)
	CreateWorkflow(spec *v1alpha1.WorkflowSpec, workflowsBatch *common.WorkflowsBatch) (*v1alpha1.Workflow, error)
	SetConfig(wf *v1alpha1.Workflow, spec *v1alpha1.WorkflowSpec) error
	Lint(wf *v1alpha1.Workflow) error
	Submit(ctx *context.Context, wf *v1alpha1.Workflow) error
	HandleWorkflowBatch(ctx *context.Context, workflowsBatch *common.WorkflowsBatch) error
}
