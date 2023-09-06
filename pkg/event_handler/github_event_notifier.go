package event_handler

import (
	"context"
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/utils"
)

var workflowTranslationToGithubMap = map[string]string{
	"":          "pending",
	"Pending":   "pending",
	"Running":   "pending",
	"Succeeded": "success",
	"Failed":    "failure",
	"Error":     "error",
}

var workflowTranslationToBitbucketMap = map[string]string{
	"":          "INPROGRESS",
	"Pending":   "INPROGRESS",
	"Running":   "INPROGRESS",
	"Succeeded": "SUCCESSFUL",
	"Failed":    "FAILED",
	"Error":     "STOPPED",
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

func (gn *githubNotifier) Notify(ctx *context.Context, workflow *v1alpha1.Workflow) error {
	fmt.Printf("Notifing workflow, %s\n", workflow.GetName())

	repo, ok := workflow.GetLabels()["repo"]
	if !ok {
		return fmt.Errorf("failed get repo label for workflow: %s", workflow.GetName())
	}
	commit, ok := workflow.GetLabels()["commit"]
	if !ok {
		return fmt.Errorf("failed get commit label for workflow: %s", workflow.GetName())
	}

	workflowLink := fmt.Sprintf("%s/workflows/%s/%s", gn.cfg.WorkflowServerConfig.ArgoAddress, gn.cfg.Namespace, workflow.GetName())

	status, err := gn.translateWorkflowStatus(string(workflow.Status.Phase), workflow.GetName())
	if err != nil {
		return err
	}

	message := utils.TrimString(workflow.Status.Message, 140) // Max length of message is 140 characters
	err = gn.clients.GitProvider.SetStatus(ctx, &repo, &commit, &workflowLink, &status, &message)
	if err != nil {
		return fmt.Errorf("failed to set status for workflow %s: %s", workflow.GetName(), err)
	}

	return nil
}

func (gn *githubNotifier) translateWorkflowStatus(status string, workflowName string) (string, error) {
	switch gn.cfg.GitProviderConfig.Provider {
	case "github":
		result, ok := workflowTranslationToGithubMap[status]
		if !ok {
			return "", fmt.Errorf("failed to translate workflow status to github stasuts for %s status: %s", workflowName, status)
		}
		return result, nil
	case "bitbucket":
		result, ok := workflowTranslationToBitbucketMap[status]
		if !ok {
			return "", fmt.Errorf("failed to translate workflow status to bitbucket stasuts for %s status: %s", workflowName, status)
		}
		return result, nil
	}
	return "", fmt.Errorf("failed to translate workflow status")
}
