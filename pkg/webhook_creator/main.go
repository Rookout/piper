package webhook_creator

import (
	"fmt"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/git_provider"
	"golang.org/x/net/context"
	"strings"
)

type WebhookCreatorImpl struct {
	clients *clients.Clients
	cfg     *conf.GlobalConfig
	hooks   []*git_provider.HookWithStatus
}

func NewWebhookCreator(cfg *conf.GlobalConfig, clients *clients.Clients) *WebhookCreatorImpl {
	wr := &WebhookCreatorImpl{
		clients: clients,
		cfg:     cfg,
	}
	return wr
}

func (wc *WebhookCreatorImpl) SetWebhooks() error {
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

func (wc *WebhookCreatorImpl) UnsetWebhooks() error {
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
