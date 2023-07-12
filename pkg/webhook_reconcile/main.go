package webhook_reconcile

import (
	"fmt"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/git_provider"
	"golang.org/x/net/context"
	"log"
	"time"
)

type WebhookReconcileImpl struct {
	clients     *clients.Clients
	cfg         *conf.GlobalConfig
	healthChan  chan *git_provider.HookWithStatus
	recoverChan chan *git_provider.HookWithStatus
	stopCh      chan struct{}
}

func NewWebhookReconcile(cfg *conf.GlobalConfig, clients *clients.Clients) *WebhookReconcileImpl {
	wr := &WebhookReconcileImpl{
		clients:     clients,
		cfg:         cfg,
		healthChan:  make(chan *git_provider.HookWithStatus),
		recoverChan: make(chan *git_provider.HookWithStatus),
		stopCh:      make(chan struct{}),
	}
	return wr
}

func Start(ctx context.Context, stop context.CancelFunc, cfg *conf.GlobalConfig, clients *clients.Clients) {
	wr := NewWebhookReconcile(cfg, clients)
	go wr.ServeAndListen(ctx)
}

func (wr *WebhookReconcileImpl) RecoverHook(hook *git_provider.HookWithStatus) error {
	ctx := context.Background()
	if hook.HealthStatus {
		return nil
	}
	recoveredHook, err := wr.clients.GitProvider.SetWebhook(&ctx, hook.RepoName)
	if err != nil {
		return err
	}
	hook.Hook = recoveredHook
	hook.HealthStatus = true
	return nil
}

func (wr *WebhookReconcileImpl) RunTest() error {
	ctx := context.Background()
	for _, hook := range wr.clients.GitProvider.GetHooks() {
		hook.HealthStatus = false
		err := wr.clients.GitProvider.PingHook(&ctx, *hook)
		if err != nil {
			log.Printf("[webhook tests] sending %v to recoveryChan", hook)
			wr.recoverChan <- hook
		}
	}
	time.Sleep(5 * time.Second) // wait for results
	for _, hook := range wr.clients.GitProvider.GetHooks() {
		if !hook.HealthStatus {
			return fmt.Errorf("[webhook tests] hook %v is not healthy", hook)
		}
	}

	return nil
}

func (wr *WebhookReconcileImpl) Stop() {
	close(wr.stopCh)
}

func (wr *WebhookReconcileImpl) Healthy(hookID *int64) error {
	for _, hook := range wr.clients.GitProvider.GetHooks() {
		if *hook.Hook.ID == *hookID {
			hook.HealthStatus = true
			return nil
		}
	}
	return fmt.Errorf("hook not found for hookdID %d", hookID)
}

func (wr *WebhookReconcileImpl) ServeAndListen(ctx context.Context) {
	defer close(wr.healthChan)
	defer close(wr.recoverChan)
	defer wr.Stop()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-wr.stopCh:
				return
			case event := <-wr.healthChan:
				if event != nil {
					log.Printf("set health status: %v", event)
					event.HealthStatus = true
				}

			case event := <-wr.recoverChan:
				if event != nil {
					log.Printf("recover health for: %v", event)
					err := wr.RecoverHook(event)
					if err != nil {
						return
					}
				}
			}
		}
	}()
}
