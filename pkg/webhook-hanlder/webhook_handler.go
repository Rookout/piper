package webhook_hanlder

import (
	"github.com/rookout/piper/pkg/conf"
)

type WebhookHandlerImpl struct {
	cfg      *conf.Config
	triggers *[]Trigger
}

func NewWebhookHandler(cfg *conf.Config, repo string, branch string) *WebhookHandlerImpl {

	return &WebhookHandlerImpl{
		cfg:      cfg,
		triggers: &[]Trigger{},
	}
}

func (wh *WebhookHandlerImpl) RegisterTriggers(cfg *conf.Config, repo string, branch string) error {
	//TODO implement me
	panic("implement me")
}

func (wh *WebhookHandlerImpl) ExecuteMatchingTriggers(event string, branch string) error {
	//TODO implement me
	panic("implement me")
}
