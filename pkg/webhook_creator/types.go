package webhook_creator

import "golang.org/x/net/context"

type WebhookCreator interface {
	Stop(ctx *context.Context)
	Start()
	SetWebhookHealth(status bool, hookID *int64) error
	RunDiagnosis(ctx *context.Context) error
}
