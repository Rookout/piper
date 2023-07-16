package common

import (
	"github.com/rookout/piper/pkg/git_provider"
)

type WorkflowsBatch struct {
	OnStart     []*git_provider.CommitFile
	OnExit      []*git_provider.CommitFile
	Templates   []*git_provider.CommitFile
	Parameters  *git_provider.CommitFile
	Config      *string
	Payload     *git_provider.WebhookPayload
	TriggerName string
}
