package webhook_reconcile

import (
	"github.com/rookout/piper/pkg/git_provider"
	"golang.org/x/net/context"
)

type WebhookReconcile interface {
	RecoverHook(hook *git_provider.HookWithStatus) error
	RunTest() error
	Stop()
	Healthy(hookId int64)
	ServeAndListen(ctx context.Context)
}
