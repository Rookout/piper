package git_provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ktrysmt/go-bitbucket"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/utils"
	"golang.org/x/net/context"
	"io"
	"log"
	"net/http"
	"strings"
)

type BitbucketServerClientImpl struct {
	client         *bitbucket.Client
	cfg            *conf.GlobalConfig
	HooksHashTable map[string]int64
}

func NewBitbucketServerClient(cfg *conf.GlobalConfig) (Client, error) {
	client := bitbucket.NewOAuthbearerToken(cfg.GitProviderConfig.Token)

	err := ValidateBitbucketPermissions(client, cfg)
	if err != nil {
		return nil, err
	}

	return &BitbucketServerClientImpl{
		client:         client,
		cfg:            cfg,
		HooksHashTable: make(map[string]int64),
	}, err
}

func (b BitbucketServerClientImpl) ListFiles(ctx *context.Context, repo string, branch string, path string) ([]string, error) {
	var filesList []string
	fileOptions := bitbucket.RepositoryFilesOptions{
		Owner:    b.cfg.GitProviderConfig.OrgName,
		RepoSlug: repo,
		Ref:      branch,
		Path:     path,
		MaxDepth: 0,
	}
	files, err := b.client.Repositories.Repository.ListFiles(&fileOptions)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		fileWithoutPath := strings.ReplaceAll(f.Path, path+"/", "")
		filesList = append(filesList, fileWithoutPath)
	}

	return filesList, nil
}

func (b BitbucketServerClientImpl) GetFile(ctx *context.Context, repo string, branch string, path string) (*CommitFile, error) {
	fileOptions := bitbucket.RepositoryFilesOptions{
		Owner:    b.cfg.GitProviderConfig.OrgName,
		RepoSlug: repo,
		Ref:      branch,
		Path:     path,
		MaxDepth: 0,
	}
	fileContent, err := b.client.Repositories.Repository.GetFileContent(&fileOptions)
	if err != nil {
		return nil, err
	}

	stringContent := string(fileContent[:])
	return &CommitFile{
		Path:    &path,
		Content: &stringContent,
	}, nil
}

func (b BitbucketServerClientImpl) GetFiles(ctx *context.Context, repo string, branch string, paths []string) ([]*CommitFile, error) {
	var commitFiles []*CommitFile
	for _, path := range paths {
		file, err := b.GetFile(ctx, repo, branch, path)
		if err != nil {
			return nil, err
		}
		if file == nil {
			log.Printf("file %s not found in repo %s branch %s", path, repo, branch)
			continue
		}
		commitFiles = append(commitFiles, file)
	}
	return commitFiles, nil
}

func (b BitbucketServerClientImpl) SetWebhook(ctx *context.Context, repo *string) (*HookWithStatus, error) {
	webhookOptions := &bitbucket.WebhooksOptions{
		Owner:       b.cfg.GitProviderConfig.OrgName,
		RepoSlug:    *repo,
		Uuid:        "",
		Description: "Piper",
		Url:         b.cfg.GitProviderConfig.WebhookURL,
		Active:      true,
		Events:      []string{"repo:push", "pullrequest:created", "pullrequest:updated", "pullrequest:fulfilled"},
	}

	hook, exists := b.isRepoWebhookExists(*repo)
	if exists {
		log.Printf("webhook already exists for repository %s, skipping creation... \n", *repo)
		addHookToHashTable(utils.RemoveBraces(hook.Uuid), b.HooksHashTable)
		hookID, err := getHookByUUID(utils.RemoveBraces(hook.Uuid), b.HooksHashTable)
		if err != nil {
			return nil, err
		}
		return &HookWithStatus{
			HookID:       hookID,
			HealthStatus: true,
			RepoName:     repo,
		}, nil
	}

	hook, err := b.client.Repositories.Webhooks.Create(webhookOptions)
	if err != nil {
		return nil, err
	}
	log.Printf("created webhook for repository %s \n", *repo)

	addHookToHashTable(utils.RemoveBraces(hook.Uuid), b.HooksHashTable)
	hookID, err := getHookByUUID(utils.RemoveBraces(hook.Uuid), b.HooksHashTable)
	if err != nil {
		return nil, err
	}

	return &HookWithStatus{
		HookID:       hookID,
		HealthStatus: true,
		RepoName:     repo,
	}, nil
}

func (b BitbucketServerClientImpl) UnsetWebhook(ctx *context.Context, hook *HookWithStatus) error {
	//TODO implement me
	panic("implement me")
}

func (b BitbucketServerClientImpl) HandlePayload(ctx *context.Context, request *http.Request, secret []byte) (*WebhookPayload, error) {
	var webhookPayload *WebhookPayload

	var buf bytes.Buffer
	_, err := io.Copy(&buf, request.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %s", err)
	}

	var body map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &body)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %s", err)
	}

	hookID, err := getHookByUUID(request.Header.Get("X-Hook-UUID"), b.HooksHashTable)
	if err != nil {
		return nil, fmt.Errorf("failed to get hook by UUID, %s", err)
	}

	// https://support.atlassian.com/bitbucket-cloud/docs/event-payloads
	switch request.Header.Get("X-Event-Key") {
	case "repo:push":
		webhookPayload = &WebhookPayload{
			Event:     "push",
			Repo:      body["repository"].(map[string]interface{})["name"].(string),
			Branch:    body["push"].(map[string]interface{})["changes"].([]interface{})[0].(map[string]interface{})["new"].(map[string]interface{})["name"].(string),
			Commit:    body["push"].(map[string]interface{})["changes"].([]interface{})[0].(map[string]interface{})["commits"].([]interface{})[0].(map[string]interface{})["hash"].(string),
			UserEmail: utils.ExtractStringsBetweenTags(body["push"].(map[string]interface{})["changes"].([]interface{})[0].(map[string]interface{})["commits"].([]interface{})[0].(map[string]interface{})["author"].(map[string]interface{})["raw"].(string))[0],
			User:      body["actor"].(map[string]interface{})["display_name"].(string),
			HookID:    hookID,
		}
	case "pullrequest:created", "pullrequest:updated":
		webhookPayload = &WebhookPayload{
			Event:            "pull_request",
			Repo:             body["repository"].(map[string]interface{})["name"].(string),
			Branch:           body["pullrequest"].(map[string]interface{})["source"].(map[string]interface{})["branch"].(map[string]interface{})["name"].(string),
			Commit:           body["pullrequest"].(map[string]interface{})["source"].(map[string]interface{})["commit"].(map[string]interface{})["hash"].(string),
			User:             body["pullrequest"].(map[string]interface{})["author"].(map[string]interface{})["display_name"].(string),
			PullRequestURL:   body["pullrequest"].(map[string]interface{})["links"].(map[string]interface{})["html"].(map[string]interface{})["href"].(string),
			PullRequestTitle: body["pullrequest"].(map[string]interface{})["title"].(string),
			DestBranch:       body["pullrequest"].(map[string]interface{})["destination"].(map[string]interface{})["branch"].(map[string]interface{})["name"].(string),
			HookID:           hookID,
		}
	}
	return webhookPayload, nil
}

func (b BitbucketServerClientImpl) SetStatus(ctx *context.Context, repo *string, commit *string, linkURL *string, status *string, message *string) error {
	commitOptions := bitbucket.CommitsOptions{
		Owner:    b.cfg.GitProviderConfig.OrgName,
		RepoSlug: *repo,
		Revision: *commit,
	}
	commitStatusOptions := bitbucket.CommitStatusOptions{
		Key:         "build",
		Url:         *linkURL,
		State:       *status,
		Description: *message,
	}
	resp, err := b.client.Repositories.Commits.CreateCommitStatus(&commitOptions, &commitStatusOptions)
	log.Printf("%s", resp)
	if err != nil {
		return err
	}

	return nil
}

func (b BitbucketServerClientImpl) PingHook(ctx *context.Context, hook *HookWithStatus) error {
	//TODO implement me
	panic("implement me")
}

func (b BitbucketServerClientImpl) isRepoWebhookExists(repo string) (*bitbucket.Webhook, bool) {
	emptyHook := bitbucket.Webhook{}

	webhookOptions := bitbucket.WebhooksOptions{
		Owner:    b.cfg.GitProviderConfig.OrgName,
		RepoSlug: repo,
	}
	hooks, err := b.client.Repositories.Webhooks.List(&webhookOptions)

	if err != nil {
		log.Printf("failed to list existing hooks for repository %s. error:%s", repo, err)
		return &emptyHook, false
	}

	if len(hooks) == 0 {
		return &emptyHook, false
	}

	for _, hook := range hooks {
		if hook.Description == "Piper" && hook.Url == b.cfg.GitProviderConfig.WebhookURL {
			return &hook, true
		}
	}

	return &emptyHook, false
}