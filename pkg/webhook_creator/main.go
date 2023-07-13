package webhook_creator

import (
	"fmt"
	"github.com/emicklei/go-restful/v3/log"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/git_provider"
	"golang.org/x/net/context"
	"strconv"
	"strings"
	"time"
)

type WebhookCreatorImpl struct {
	clients *clients.Clients
	cfg     *conf.GlobalConfig
	hooks   map[int64]*git_provider.HookWithStatus
}

func NewWebhookCreator(cfg *conf.GlobalConfig, clients *clients.Clients) *WebhookCreatorImpl {
	wr := &WebhookCreatorImpl{
		clients: clients,
		cfg:     cfg,
		hooks:   make(map[int64]*git_provider.HookWithStatus, 0),
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
	hook, ok := wc.hooks[*hookID]
	if !ok {
		return fmt.Errorf("unable to set health status for hookID %d", *hookID)
	}
	hook.HealthStatus = status
	log.Printf("set health status to %s for hook id: %d", strconv.FormatBool(status), *hookID)
	return nil
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
		wc.hooks[*hook.HookID] = hook
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
	log.Printf("Starting webhook diagnostics")
	wc.setAllHooksHealth(false)
	wc.pingHooks(ctx)
	wc.checkHooksHealth(10 * time.Second)
	for hookID, hook := range wc.hooks {
		if !hook.HealthStatus {
			log.Printf("Trying to recover hook %d", hookID)
			err := wc.recoverHook(ctx, hook.HookID)
			if err != nil {
				return err
			}
		}
	}
	log.Print("Successful webhook diagnosis")
	return nil
}

func (wc *WebhookCreatorImpl) pingHooks(ctx *context.Context) {
	for hookID, hook := range wc.hooks {
		err := wc.clients.GitProvider.PingHook(ctx, hook)
		if err != nil {
			log.Printf("failed to ping hook: %v", err)
			log.Printf("Trying to recover from ping hook %d", hookID)
			err = wc.recoverHook(ctx, &hookID)
			if err != nil {
				log.Printf("failed recover hookID:%d got error:%s", hookID, err)
			}
		}
	}
}
func (wc *WebhookCreatorImpl) setAllHooksHealth(status bool) {
	for _, hook := range wc.hooks {
		hook.HealthStatus = status
	}
	log.Printf("set all hooks health status for to %s", strconv.FormatBool(status))
}

func (wc *WebhookCreatorImpl) checkHooksHealth(timeout time.Duration) bool {
	startTime := time.Now()

	for {
		allHealthy := true
		for _, hook := range wc.hooks {
			if !hook.HealthStatus {
				allHealthy = false
				break
			}
		}

		if allHealthy {
			return true
		}

		if time.Since(startTime) >= timeout {
			break
		}

		time.Sleep(1 * time.Second) // Adjust the sleep duration as per your requirement
	}

	return false
}
