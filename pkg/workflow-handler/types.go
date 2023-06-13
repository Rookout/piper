package workflow_handler

import "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"

type Client interface {
	NewTemplate() (*v1alpha1.Template, error)
	NewSpec(templates []*v1alpha1.Template) (*v1alpha1.WorkflowSpec, error)
	NewWorkflow(spec *v1alpha1.WorkflowSpec) (*v1alpha1.Workflow, error)
	SetConfig(wf *v1alpha1.Workflows, spec *v1alpha1.WorkflowSpec) error
	Lint(wf *v1alpha1.Workflows) error
	Submit(wf *v1alpha1.Workflows) error
	//GetWorkflowStatus(wf *v1alpha1.Workflows) (*v1alpha1.WorkflowStatus, error)
}
