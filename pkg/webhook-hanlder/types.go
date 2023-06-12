package webhook_hanlder

type WebhookHandler interface {
	RegisterTriggers() error
	ExecuteMatchingTriggers() error
}
