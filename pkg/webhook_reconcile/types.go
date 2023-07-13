package webhook_reconcile

import (
	"golang.org/x/net/context"
)

type WebhookReconcile interface {
	RecoverHook(hookID *int64) error
	Healthy(hookId int64)
	RunTest() error
	Stop()
	ServeAndListen(ctx context.Context)
}
