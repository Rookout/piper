package workflow_handler

import (
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	wfClientSet "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/utils"
)

type WorkflowsClientImpl struct {
	clientSet *wfClientSet.Clientset
	cfg       *conf.Config
}

func NewWorkflowsClient(cfg *conf.Config) (*WorkflowsClientImpl, error) {
	restClientConfig, err := utils.GetClientConfig(cfg.ArgoConfig.KubeConfig)
	if err != nil {
		return nil, err
	}

	clientSet := wfClientSet.NewForConfigOrDie(restClientConfig) //.ArgoprojV1alpha1().Workflows(namespace)
	return &WorkflowsClientImpl{
		clientSet: clientSet,
		cfg:       cfg,
	}, nil
}

func (wfc *WorkflowsClientImpl) NewTemplate() (*v1alpha1.Template, error) {
	//TODO implement me
	panic("implement me")
}

func (wfc *WorkflowsClientImpl) NewSpec(templates []*v1alpha1.Template) (*v1alpha1.WorkflowSpec, error) {
	//TODO implement me
	panic("implement me")
}

func (wfc *WorkflowsClientImpl) NewWorkflow(sepc *v1alpha1.WorkflowSpec) (*v1alpha1.Workflow, error) {
	//TODO implement me
	panic("implement me")
}

func (wfc *WorkflowsClientImpl) SetConfig(wf *v1alpha1.Workflows, spec *v1alpha1.WorkflowSpec) error {
	//TODO implement me
	panic("implement me")
}

func (wfc *WorkflowsClientImpl) Lint(wf *v1alpha1.Workflows) error {
	//TODO implement me
	panic("implement me")
}

func (wfc *WorkflowsClientImpl) Submit(wf *v1alpha1.Workflows) error {
	//TODO implement me
	panic("implement me")
}
