package webhook_hanlder

import (
	"github.com/rookout/piper/pkg/conf"
)

type Trigger struct {
	events   []string `json:"events"`
	branches []string `json:"branches"`
	onStart  []string `json:"execute"`
	onExit   []string `json:"on_exit"`
}

type WebhookHandler interface {
	RegisterTriggers(cfg *conf.Config, repo string, branch string) error
	ExecuteMatchingTriggers(event string, branch string) error
}
