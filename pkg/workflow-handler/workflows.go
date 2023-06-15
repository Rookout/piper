package workflow_handler

import (
	"context"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	wfClientSet "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"

	"github.com/rookout/piper/pkg/common"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/utils"
)

const (
	ENTRYPOINT = "entryPoint"
	ONEXIT     = "exitHandler"
)

type WorkflowsClientImpl struct {
	clientSet *wfClientSet.Clientset
	cfg       *conf.Config
}

func NewWorkflowsClient(cfg *conf.Config) (Client, error) {
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

func (wfc *WorkflowsClientImpl) CreateTemplate(workflowsBatch *WorkflowsBatch) ([]v1alpha1.Template, error) {
	finalTemplate := make([]v1alpha1.Template, 0)
	onStart, err := CreateDAGTemplate(workflowsBatch.OnStart, ENTRYPOINT)
	if err != nil {
		return nil, err
	}
	finalTemplate = append(finalTemplate, *onStart)

	onExit, err := CreateDAGTemplate(workflowsBatch.OnExit, ONEXIT)
	if err != nil {
		return nil, err
	}
	finalTemplate = append(finalTemplate, *onExit)

	finalTemplate, err = AddFilesToTemplate(finalTemplate, workflowsBatch.Templates)
	if err != nil {
		return nil, err
	}

	return finalTemplate, nil
}

func (wfc *WorkflowsClientImpl) CreateSpec(templates []v1alpha1.Template) (*v1alpha1.WorkflowSpec, error) {
	finalSpec := &v1alpha1.WorkflowSpec{}
	//err := yaml.Unmarshal([]byte(*f.Content), finalSpec)
	//if err != nil {
	//	return nil, err
	//}

	return finalSpec, nil
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

func HandleWorkflowBatch(ctx *context.Context, wfc Client, workflowsBatch *common.WorkflowsBatch) error {

	return nil
}
