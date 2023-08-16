package git_provider

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	bitbucket "github.com/gfleury/go-bitbucket-v1"
	"github.com/mitchellh/mapstructure"
	"github.com/rookout/piper/pkg/conf"
	"golang.org/x/net/context"
	"io"
	"log"
	"net/http"
)

type BitbucketServerClientImpl struct {
	client *bitbucket.APIClient
	cfg    *conf.GlobalConfig
}

func NewBitbucketServerClient(cfg *conf.GlobalConfig) (Client, error) {
	bitbucketConfig := &bitbucket.Configuration{
		BasePath:      cfg.GitProviderConfig.BaseURL,
		DefaultHeader: make(map[string]string),
		UserAgent:     "piper",
	}
	bitbucketConfig.AddDefaultHeader("x-atlassian-token", "no-check")
	bitbucketConfig.AddDefaultHeader("x-requested-with", "XMLHttpRequest")

	ctx := context.Background() // context.WithTimeout(context.Background(), 10*time.Second)

	ctx = context.WithValue(ctx, bitbucket.ContextAccessToken, cfg.GitProviderConfig.Token)
	client := bitbucket.NewAPIClient(ctx, bitbucketConfig)

	err := BitbucketValidatePermissions(ctx, client, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to validate permissions: %v", err)
	}

	return &BitbucketServerClientImpl{
		client: client,
		cfg:    cfg,
	}, err
}

func (b BitbucketServerClientImpl) ListFiles(ctx *context.Context, repo string, branch string, path string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (b BitbucketServerClientImpl) GetFile(ctx *context.Context, repo string, branch string, path string) (*CommitFile, error) {
	//TODO implement me
	panic("implement me")
}

func (b BitbucketServerClientImpl) GetFiles(ctx *context.Context, repo string, branch string, paths []string) ([]*CommitFile, error) {
	//TODO implement me
	panic("implement me")
}

func (b BitbucketServerClientImpl) SetWebhook(ctx *context.Context, repo *string) (*HookWithStatus, error) {
	if b.cfg.OrgLevelWebhook && repo != nil {
		return nil, fmt.Errorf("trying to set repo scope. repo: %s", *repo)
	}

	if repo == nil {
		return nil, fmt.Errorf("org scope not supported")
	} else {
		newHook := bitbucket.Webhook{
			Name:   "Piper",
			Url:    b.cfg.WebhookURL,
			Active: true,
			// https://confluence.atlassian.com/bitbucketserver/event-payload-938025882.html
			Events:        []string{"repo:refs_changed", "pr:opened", "pr:merged"},
			Configuration: bitbucket.WebhookConfiguration{Secret: b.cfg.GitProviderConfig.WebhookSecret},
		}
		requestBody, err := json.Marshal(newHook)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal new webhook to JSON, %w", err)
		}

		respHook, ok := b.isRepoWebhookExists(*ctx, *repo)
		if !ok {

			resp, err := b.client.DefaultApi.CreateWebhook(b.cfg.GitProviderConfig.OrgName, *repo, requestBody, []string{"application/json"})
			if err != nil {
				return nil, fmt.Errorf("failed to add webhook. error: %w", err)
			}

			var createdHook *bitbucket.Webhook
			err = mapstructure.Decode(resp.Values, &createdHook)
			if err != nil {
				return nil, fmt.Errorf("failed to convert API response to Webhook struct. error: %w", err)
			}

			log.Printf("created webhook  for %s: %s\n", *repo, createdHook.Url)

			return &HookWithStatus{HookID: int64(createdHook.ID), HealthStatus: true, RepoName: repo}, nil
		} else {
			if respHook.Configuration.Secret != b.cfg.GitProviderConfig.WebhookSecret {
				resp, err := b.client.DefaultApi.UpdateWebhook(b.cfg.GitProviderConfig.OrgName, *repo, int32(respHook.ID), requestBody, []string{"application/json"})
				if err != nil {
					return nil, fmt.Errorf("failed to add update webhook. error: %w", err)
				}

				var updatedHook *bitbucket.Webhook
				err = mapstructure.Decode(resp.Values, &updatedHook)
				if err != nil {
					return nil, fmt.Errorf("failed to convert API response to Webhook struct. error: %w", err)
				}

				log.Printf("updated webhook  for %s: %s\n", *repo, updatedHook.Url)
				return &HookWithStatus{HookID: int64(updatedHook.ID), HealthStatus: true, RepoName: repo}, nil
			}
			log.Printf("webhook exists for %s: %s\n", *repo, respHook.Url)
			return &HookWithStatus{HookID: int64(respHook.ID), HealthStatus: true, RepoName: repo}, nil
		}
	}
}

func (b BitbucketServerClientImpl) UnsetWebhook(ctx *context.Context, hook *HookWithStatus) error {
	//TODO implement me
	panic("implement me")
}

func (b BitbucketServerClientImpl) HandlePayload(ctx *context.Context, request *http.Request, secret []byte) (*WebhookPayload, error) {
	var webhookPayload *WebhookPayload

	eventType := request.Header.Get("X-Event-Key")
	if eventType == "diagnostics:ping" {
		return &WebhookPayload{
			Event: "ping",
		}, nil
	}
	_, err := b.validatRequest(request, secret)
	if err != nil {
		return nil, err
	}

	return webhookPayload, nil
}

func (b BitbucketServerClientImpl) SetStatus(ctx *context.Context, repo *string, commit *string, linkURL *string, status *string, message *string) error {
	//TODO implement me
	panic("implement me")
}

func (b BitbucketServerClientImpl) PingHook(ctx *context.Context, hook *HookWithStatus) error {

	body := map[string]interface{}{
		"url": b.cfg.GitProviderConfig.WebhookURL,
	}

	apiResponse, err := b.client.DefaultApi.TestWebhook(b.cfg.GitProviderConfig.OrgName, *hook.RepoName, body)
	if err != nil {
		return err
	}

	if apiResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook of repo %s test returned %s", *hook.RepoName, apiResponse.Message)
	}

	return nil
}

func (b BitbucketServerClientImpl) isRepoWebhookExists(ctx context.Context, repo string) (*bitbucket.Webhook, bool) {
	emptyHook := bitbucket.Webhook{}
	apiResponse, err := b.client.DefaultApi.FindWebhooks(b.cfg.GitProviderConfig.OrgName, repo, nil)
	if err != nil {
		log.Printf("failed to list existing hooks for repository %s. error:%s", repo, err)
		return &emptyHook, false
	}
	if apiResponse.StatusCode != 200 {
		return &emptyHook, false
	}

	hooks, err := bitbucket.GetWebhooksResponse(apiResponse)
	if err != nil {
		log.Printf("failed to convert the list of webhooks for repository %s. error:%s", repo, err)
		return &emptyHook, false
	}

	if len(hooks) == 0 {
		return &emptyHook, false
	}

	for _, hook := range hooks {
		if hook.Name == "Piper" && hook.Url == b.cfg.GitProviderConfig.WebhookURL {
			return &hook, true
		}
	}

	return &emptyHook, false
}

func (b *BitbucketServerClientImpl) validatRequest(request *http.Request, secret []byte) ([]byte, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request body, %w", err)
	}

	signature := request.Header.Get("X-Hub-Signature")
	if len(signature) == 0 {
		return nil, fmt.Errorf("missing signature header")
	}

	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write(body)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(signature[7:]), []byte(expectedMAC)) {
		return nil, fmt.Errorf("hmac verification failed")
	}

	return body, nil
}
