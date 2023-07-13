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
	clients *clients.Clients
	cfg     *conf.GlobalConfig
	hooks   []*git_provider.HookWithStatus
}

func NewWebhookCreator(cfg *conf.GlobalConfig, clients *clients.Clients) *WebhookCreatorImpl {
	wr := &WebhookCreatorImpl{
		clients: clients,
		cfg:     cfg,
	}

	err := wr.setWebhooks()
	if err != nil {
		log.Print(err)
		panic("failed in initializing webhooks")
	}
	return wr
}

func (wc *WebhookCreatorImpl) recoverHook(ctx *context.Context, hookID *int64) error {
	for i, hook := range wc.hooks {
		if *hook.HookID == *hookID {
			newHook, err := wc.clients.GitProvider.SetWebhook(ctx, hook.RepoName)
			if err != nil {
				return err
			}
			wc.hooks[i] = newHook
			return nil
		}
	}
	return fmt.Errorf("unable to recover hookID %d, not found in list of hooks", hookID)
}

func (wc *WebhookCreatorImpl) SetHealth(status bool, hookID *int64) error {
	for _, hook := range wc.hooks {
		if *hook.HookID == *hookID {
			hook.HealthStatus = status
			log.Printf("set health status to %b for hook id: %d", status, *hookID)
			return nil
		}
	}
	return fmt.Errorf("unable to set health status for hookID %d", *hookID)
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

func (wc *WebhookCreatorImpl) unsetWebhooks(ctx *context.Context) error {
	for _, hook := range wc.hooks {
		err := wc.clients.GitProvider.UnsetWebhook(ctx, hook)
		if err != nil {
			return err
		}
	}

	return nil
}

func (wc *WebhookCreatorImpl) Stop(ctx *context.Context) {
	err := wc.unsetWebhooks(ctx)
	if err != nil {
		log.Printf("Failed to unset webhooks, error: %v", err)
	}
}

func (wc *WebhookCreatorImpl) RunDiagnosis(ctx *context.Context) error {
	wc.setAllHooksHealth(false)
	wc.pingHooks(ctx)
	for _, hook := range wc.hooks {
		if !hook.HealthStatus {
			log.Printf("Trying to recover hook %d", hook.HookID)
			err := wc.recoverHook(ctx, hook.HookID)
			if err != nil {
				return err
			}
			return fmt.Errorf("failed webhook diagnosis: hook %d is not healthy", hook.HookID)
		}
	}
	log.Print("Successful webhook diagnosis")
	return nil
}

func (wc *WebhookCreatorImpl) pingHooks(ctx *context.Context) {
	for _, hook := range wc.hooks {
		err := wc.clients.GitProvider.PingHook(ctx, hook)
		if err != nil {
			log.Printf("failed to ping hook: %v", err)
			log.Printf("Trying to recover hook %d", hook.HookID)
			err = wc.recoverHook(ctx, hook.HookID)
			if err != nil {
				log.Printf("failed recover hookID:%d got error:%s", hook.HookID, err)
			}
		}
	}
}
func (wc *WebhookCreatorImpl) setAllHooksHealth(status bool) {
	for _, hook := range wc.hooks {
		hook.HealthStatus = status
	}
	log.Printf("set all hooks health status for to %b", status)
}
