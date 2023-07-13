package webhook_creator

type WebhookCreator interface {
	SetWebhooks() error
	UnsetWebhooks() error
	Shutdown()
}
