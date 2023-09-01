package git_provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ktrysmt/go-bitbucket"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/utils"
	"github.com/tidwall/gjson"
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
			Repo:      gjson.GetBytes(buf.Bytes(), "repository.name").Value().(string),
			Branch:    gjson.GetBytes(buf.Bytes(), "push.changes.0.new.name").Value().(string),
			Commit:    gjson.GetBytes(buf.Bytes(), "push.changes.0.commits.0.hash").Value().(string),
			UserEmail: utils.ExtractStringsBetweenTags(gjson.GetBytes(buf.Bytes(), "push.changes.0.commits.0.author.raw").Value().(string))[0],
			User:      gjson.GetBytes(buf.Bytes(), "actor.display_name").Value().(string),
			HookID:    hookID,
		}
	case "pullrequest:created", "pullrequest:updated":
		webhookPayload = &WebhookPayload{
			Event:            "pull_request",
			Repo:             gjson.GetBytes(buf.Bytes(), "repository.name").Value().(string),
			Branch:           gjson.GetBytes(buf.Bytes(), "pullrequest.source.branch.name").Value().(string),
			Commit:           gjson.GetBytes(buf.Bytes(), "pullrequest.source.commit.hash").Value().(string),
			User:             gjson.GetBytes(buf.Bytes(), "pullrequest.author.display_name").Value().(string),
			PullRequestURL:   gjson.GetBytes(buf.Bytes(), "pullrequest.links.html.href").Value().(string),
			PullRequestTitle: gjson.GetBytes(buf.Bytes(), "pullrequest.title").Value().(string),
			DestBranch:       gjson.GetBytes(buf.Bytes(), "pullrequest.destination.branch.name").Value().(string),
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
	_, err := b.client.Repositories.Commits.CreateCommitStatus(&commitOptions, &commitStatusOptions)
	if err != nil {
		return err
	}
	log.Printf("set status of commit %s in repo %s to %s", *commit, *repo, *status)
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
