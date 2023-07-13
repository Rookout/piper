package webhook_creator

import "golang.org/x/net/context"

type WebhookCreator interface {
	Stop(ctx *context.Context)
	Start()
	SetToHealthy(hookID *int64) error
}
