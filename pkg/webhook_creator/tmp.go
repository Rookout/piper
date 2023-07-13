package webhook_creator

//
//import (
//	"fmt"
//	"github.com/rookout/piper/pkg/clients"
//	"github.com/rookout/piper/pkg/conf"
//	"github.com/rookout/piper/pkg/git_provider"
//	"golang.org/x/net/context"
//	"log"
//	"strings"
//)
//
//type WebhookReconcileImpl struct {
//	clients           *clients.Clients
//	cfg               *conf.GlobalConfig
//	hookIDHealthChan  chan *int64
//	hookIDRecoverChan chan *int64
//	stopCh            chan struct{}
//}
//
//func NewWebhookReconcile(cfg *conf.GlobalConfig, clients *clients.Clients) *WebhookReconcileImpl {
//	wr := &WebhookReconcileImpl{
//		clients:           clients,
//		cfg:               cfg,
//		hookIDHealthChan:  make(chan *int64),
//		hookIDRecoverChan: make(chan *int64),
//		stopCh:            make(chan struct{}),
//	}
//	return wr
//}
//
//func (wr *WebhookReconcileImpl) RecoverHook(hookID *int64) error {
//	ctx := context.Background()
//	hook, err := wr.getHook(hookID)
//	if err != nil {
//		return err
//	}
//	if hook.HealthStatus {
//		return nil
//	}
//	recoveredHook, err := wr.clients.GitProvider.SetWebhook(&ctx, hook.RepoName)
//	if err != nil {
//		return err
//	}
//	hook.HookID = recoveredHook.HookID
//	hook.HealthStatus = true
//	return nil
//}
//
//func (wr *WebhookReconcileImpl) RunTest() error {
//	ctx := context.Background()
//	for _, hook := range wr.clients.GitProvider.GetHooks() {
//		hook.HealthStatus = false
//		err := wr.clients.GitProvider.PingHook(&ctx, *hook)
//		if err != nil {
//			log.Printf("[webhook tests] error: %s", err)
//			log.Printf("[webhook tests] sending %v to recoveryChan", hook)
//			wr.hookIDRecoverChan <- hook.HookID
//		}
//	}
//
//	for _, hook := range wr.clients.GitProvider.GetHooks() {
//		if !hook.HealthStatus {
//			return fmt.Errorf("[webhook tests] hook %v is not healthy", hook)
//		}
//	}
//	return nil
//}
//
//func (wr *WebhookReconcileImpl) Stop() {
//	close(wr.stopCh)
//}
//
//func (wr *WebhookReconcileImpl) Healthy(hookID *int64) error {
//	hook, err := wr.getHook(hookID)
//	if err != nil {
//		return err
//	}
//
//	hook.HealthStatus = true
//	return nil
//}
//
//func (wr *WebhookReconcileImpl) getHook(hookID *int64) (*git_provider.HookWithStatus, error) {
//	for _, hook := range wr.clients.GitProvider.GetHooks() {
//		if *hook.HookID == *hookID {
//			return hook, nil
//		}
//	}
//	return nil, fmt.Errorf("hook with hoodID:%d not found", *hookID)
//}
//

//
//func (wr *WebhookReconcileImpl) ServeAndListen(ctx context.Context) {
//	defer close(wr.hookIDHealthChan)
//	defer close(wr.hookIDRecoverChan)
//	defer wr.Stop()
//	go func() {
//		for {
//			select {
//			case <-ctx.Done():
//				return
//			case <-wr.stopCh:
//				return
//			case hookID := <-wr.hookIDHealthChan:
//				if hookID != nil {
//					log.Printf("set health status for hook id: %d", hookID)
//					wr.Healthy(hookID)
//				}
//			case hookID := <-wr.hookIDRecoverChan:
//				if hookID != nil {
//					log.Printf("recover health for hook id: %d", hookID)
//					err := wr.RecoverHook(hookID)
//					if err != nil {
//						return
//					}
//				}
//			}
//		}
//	}()
//}
