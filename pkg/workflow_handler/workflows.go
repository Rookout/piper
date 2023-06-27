package workflow_handler

import (
	"context"
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	wfClientSet "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"

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
	cfg       *conf.GlobalConfig
}

func NewWorkflowsClient(cfg *conf.GlobalConfig) (WorkflowsClient, error) {
	restClientConfig, err := utils.GetClientConfig(cfg.WorkflowServerConfig.KubeConfig)
	if err != nil {
		return nil, err
	}

	clientSet := wfClientSet.NewForConfigOrDie(restClientConfig) //.ArgoprojV1alpha1().Workflows(namespace)
	return &WorkflowsClientImpl{
		clientSet: clientSet,
		cfg:       cfg,
	}, nil
}

func (wfc *WorkflowsClientImpl) ConstructTemplates(workflowsBatch *common.WorkflowsBatch, configName string) ([]v1alpha1.Template, error) {
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
	if onExit == nil || len(onExit.DAG.Tasks) == 0 {
		if IsConfigExists(&wfc.cfg.WorkflowsConfig, configName) && IsConfigsOnExitExists(&wfc.cfg.WorkflowsConfig, configName) {
			template := &v1alpha1.Template{
				Name: ONEXIT,
				DAG: &v1alpha1.DAGTemplate{
					Tasks: wfc.cfg.WorkflowsConfig.Configs[configName].OnExit,
				},
			}

			finalTemplate = append(finalTemplate, *template)
		}
	} else {
		finalTemplate = append(finalTemplate, *onExit)
	}

	finalTemplate, err = AddFilesToTemplates(finalTemplate, workflowsBatch.Templates)
	if err != nil {
		return nil, err
	}

	return finalTemplate, nil
}

func (wfc *WorkflowsClientImpl) ConstructSpec(templates []v1alpha1.Template, params []v1alpha1.Parameter, configName string) (*v1alpha1.WorkflowSpec, error) {
	finalSpec := &v1alpha1.WorkflowSpec{}
	if IsConfigExists(&wfc.cfg.WorkflowsConfig, configName) {
		*finalSpec = wfc.cfg.WorkflowsConfig.Configs[configName].Spec
		if len(wfc.cfg.WorkflowsConfig.Configs[configName].OnExit) != 0 {
			finalSpec.OnExit = ONEXIT
		}
	}

	finalSpec.Entrypoint = ENTRYPOINT
	finalSpec.Templates = templates
	finalSpec.Arguments.Parameters = params

	return finalSpec, nil
}

func (wfc *WorkflowsClientImpl) CreateWorkflow(spec *v1alpha1.WorkflowSpec, workflowsBatch *common.WorkflowsBatch) (*v1alpha1.Workflow, error) {
	workflow := &v1alpha1.Workflow{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: ConvertToValidString(workflowsBatch.Payload.Repo + "-" + workflowsBatch.Payload.Branch + "-"),
			Namespace:    wfc.cfg.Namespace,
			Labels: map[string]string{
				"repo":   workflowsBatch.Payload.Repo,
				"branch": workflowsBatch.Payload.Branch,
				"user":   workflowsBatch.Payload.User,
				"commit": workflowsBatch.Payload.Commit,
			},
		},
		Spec: *spec,
	}

	return workflow, nil
}

func (wfc *WorkflowsClientImpl) SelectConfig(workflowsBatch *common.WorkflowsBatch) (string, error) {
	var configName string
	if IsConfigExists(&wfc.cfg.WorkflowsConfig, "default") {
		configName = "default"
	}

	if *workflowsBatch.Config != "" {
		if IsConfigExists(&wfc.cfg.WorkflowsConfig, *workflowsBatch.Config) {
			configName = *workflowsBatch.Config
		} else {
			return configName, fmt.Errorf(
				"error in selecting config, staying with default config for repo %s branch %s",
				workflowsBatch.Payload.Repo,
				workflowsBatch.Payload.Branch,
			)
		}
	}

	log.Printf(
		"%s config selected for workflow in repo: %s branch %s",
		configName,
		workflowsBatch.Payload.Repo,
		workflowsBatch.Payload.Branch,
	) // Info

	return configName, nil
}

func (wfc *WorkflowsClientImpl) Lint(wf *v1alpha1.Workflow) error {
	//TODO implement me
	panic("implement me")
}

func (wfc *WorkflowsClientImpl) Submit(ctx *context.Context, wf *v1alpha1.Workflow) error {
	workflowsClient := wfc.clientSet.ArgoprojV1alpha1().Workflows(wfc.cfg.Namespace)
	_, err := workflowsClient.Create(*ctx, wf, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (wfc *WorkflowsClientImpl) HandleWorkflowBatch(ctx *context.Context, workflowsBatch *common.WorkflowsBatch) error {
	var params []v1alpha1.Parameter

	configName, err := wfc.SelectConfig(workflowsBatch)
	if err != nil {
		return err
	}

	templates, err := wfc.ConstructTemplates(workflowsBatch, configName)
	if err != nil {
		return err
	}

	if workflowsBatch.Parameters != nil {
		params, err = GetParameters(workflowsBatch.Parameters)
		if err != nil {
			return err
		}
	}

	globalParams := []v1alpha1.Parameter{
		{Name: "dest_branch", Value: v1alpha1.AnyStringPtr(workflowsBatch.Payload.DestBranch)},
		{Name: "commit", Value: v1alpha1.AnyStringPtr(workflowsBatch.Payload.Commit)},
		{Name: "branch", Value: v1alpha1.AnyStringPtr(workflowsBatch.Payload.Branch)},
		{Name: "repo_name", Value: v1alpha1.AnyStringPtr(workflowsBatch.Payload.Repo)},
		{Name: "event_type", Value: v1alpha1.AnyStringPtr(workflowsBatch.Payload.Event)},
		{Name: "pull_request_title", Value: v1alpha1.AnyStringPtr(workflowsBatch.Payload.PullRequestTitle)},
		{Name: "pull_request_url", Value: v1alpha1.AnyStringPtr(workflowsBatch.Payload.PullRequestURL)},
	}
	params = append(params, globalParams...)

	spec, err := wfc.ConstructSpec(templates, params, configName)
	if err != nil {
		return err
	}

	workflow, err := wfc.CreateWorkflow(spec, workflowsBatch)
	if err != nil {
		return err
	}

	err = wfc.Submit(ctx, workflow)
	if err != nil {
		return fmt.Errorf("failed to submit workflow, error: %v", err)
	}

	log.Printf("submit workflow for branch %s repo %s commit %s", workflowsBatch.Payload.Branch, workflowsBatch.Payload.Repo, workflowsBatch.Payload.Commit)
	return nil
}
