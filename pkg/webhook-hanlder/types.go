package webhook_hanlder

import (
	"context"
	"github.com/rookout/piper/pkg/common"
)

type Trigger struct {
	Events    *[]string `yaml:"events"`
	Branches  *[]string `yaml:"branches"`
	OnStart   *[]string `yaml:"onStart"`
	Templates *[]string `yaml:"templates"`
	OnExit    *[]string `yaml:"onExit"`
	Config    string    `yaml:"config" default:"default"`
}

type WebhookHandler interface {
	RegisterTriggers(ctx *context.Context) error
	PrepareBatchForMatchingTriggers(ctx *context.Context) ([]*common.WorkflowsBatch, error)
}
