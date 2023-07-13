package webhook_creator

import (
	"fmt"
	"github.com/emicklei/go-restful/v3/log"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/git_provider"
	"golang.org/x/net/context"
	"strings"
)

type WebhookCreatorImpl struct {
	clients          *clients.Clients
	cfg              *conf.GlobalConfig
	hooks            []*git_provider.HookWithStatus
	hookIDHealthChan chan *int64
	stopChan         *SafeChannel
}

func NewWebhookCreator(cfg *conf.GlobalConfig, clients *clients.Clients) *WebhookCreatorImpl {
	wr := &WebhookCreatorImpl{
		clients:          clients,
		cfg:              cfg,
		hookIDHealthChan: make(chan *int64),
		stopChan:         NewSafeChannel(),
	}

	err := wr.setWebhooks()
	if err != nil {
		log.Print(err)
		panic("failed in initializing webhooks")
	}
	return wr
}

func (wc *WebhookCreatorImpl) Start() {
	go func() {
		for {
			select {
			case <-wc.stopChan.C:
				log.Printf("Deleting webhooks")
				close(wc.hookIDHealthChan)
				err := wc.unsetWebhooks()
				if err != nil {
					log.Printf("Failed to unset webhooks, error: %v", err)
				}
				return
			case hookID := <-wc.hookIDHealthChan:
				if hookID != nil {
					log.Printf("set health status for hook id: %d", hookID)
					//wc.Healthy(hookID)
				}
			}
		}
	}()
}

func (wc *WebhookCreatorImpl) setWebhooks() error {
	ctx := context.Background()
	if wc.cfg.GitProviderConfig.OrgLevelWebhook && len(wc.cfg.GitProviderConfig.RepoList) != 0 {
		return fmt.Errorf("org level webhook wanted but provided repositories list")
	}
	for _, repo := range strings.Split(wc.cfg.GitProviderConfig.RepoList, ",") {
		hook, err := wc.clients.GitProvider.SetWebhook(&ctx, &repo)
		if err != nil {
			return err
		}
		wc.hooks = append(wc.hooks, hook)
	}

	return nil
}

func (wc *WebhookCreatorImpl) unsetWebhooks() error {
	ctx := context.Background()
	if wc.cfg.GitProviderConfig.OrgLevelWebhook && len(wc.cfg.GitProviderConfig.RepoList) != 0 {
		return fmt.Errorf("org level webhook wanted but provided repositories list")
	}
	for _, repo := range strings.Split(wc.cfg.GitProviderConfig.RepoList, ",") {
		hook, err := wc.clients.GitProvider.SetWebhook(&ctx, &repo)
		if err != nil {
			return err
		}
		wc.hooks = append(wc.hooks, hook)
	}

	return nil
}

func (wc *WebhookCreatorImpl) Stop() {
	wc.stopChan.C <- struct{}{}
}
