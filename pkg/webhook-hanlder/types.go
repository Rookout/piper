package webhook_hanlder

import (
	"context"
	"github.com/rookout/piper/pkg/git"
)

type WorkflowsBatch struct {
	OnStart    []*git.CommitFile
	OnExit     []*git.CommitFile
	Parameters *git.CommitFile
}

type Trigger struct {
	Events   *[]string `yaml:"events"`
	Branches *[]string `yaml:"branches"`
	OnStart  *[]string `yaml:"onStart"`
	OnExit   *[]string `yaml:"onExit"`
}

type WebhookHandler interface {
	RegisterTriggers(ctx *context.Context) error
	PrepareBatchForMatchingTriggers(ctx *context.Context) ([]*WorkflowsBatch, error)
}
