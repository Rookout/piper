package event_handler

import (
	"context"
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
)

var workflowTranslationToGithubMap = map[string]string{
	"":          "error",
	"Pending":   "pending",
	"Running":   "pending",
	"Succeeded": "success",
	"Failed":    "failure",
	"Error":     "error",
}

type githubNotifier struct {
	cfg     *conf.GlobalConfig
	clients *clients.Clients
}

func NewGithubEventNotifier(cfg *conf.GlobalConfig, clients *clients.Clients) EventNotifier {
	return &githubNotifier{
		cfg:     cfg,
		clients: clients,
	}
}

func (gn *githubNotifier) notify(ctx *context.Context, workflow *v1alpha1.Workflow) error {
	fmt.Printf("Notifing workflow, %s\n", workflow.GetName())

	repo, ok := workflow.GetLabels()["repo"]
	if !ok {
		return fmt.Errorf("failed get repo label for workflow: %s", workflow.GetName())
	}
	commit, ok := workflow.GetLabels()["commit"]
	if !ok {
		return fmt.Errorf("failed get commit label for workflow: %s", workflow.GetName())
	}

	// TODO: separate internal and external workflow addresses
	workflowLink := fmt.Sprintf("%s/workflows/%s/%s", gn.cfg.WorkflowServerConfig.ArgoAddress, gn.cfg.Namespace, workflow.GetName())

	status, ok := workflowTranslationToGithubMap[string(workflow.Status.Phase)]
	if !ok {
		return fmt.Errorf("failed to translate workflow status to github stasuts for %s status: %s", workflow.GetName(), workflow.Status.Phase)
	}

	message := workflow.Status.Message
	err := gn.clients.GitProvider.SetStatus(ctx, &repo, &commit, &workflowLink, &status, &message)
	if err != nil {
		return fmt.Errorf("failed to set status for workflow %s: %s", workflow.GetName(), err)
	}

	return nil
}
