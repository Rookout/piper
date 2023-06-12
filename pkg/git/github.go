package git

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/rookout/piper/pkg/conf"

	"github.com/google/go-github/v52/github"
)

type GithubClientImpl struct {
	client *github.Client
	cfg    *conf.Config
	hooks  []*github.Hook
}

func NewGithubClient(cfg *conf.Config) (Client, error) {
	ctx := context.Background()

	client := github.NewTokenClient(ctx, cfg.GitConfig.Token)
	err := ValidatePermissions(ctx, client, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to validate permissions: %v", err)
	}

	return &GithubClientImpl{
		client: client,
		cfg:    cfg,
		hooks:  []*github.Hook{},
	}, err
}

func (c *GithubClientImpl) ListFiles(repo string, branch string, path string) ([]string, error) {
	var files []string
	ctx := context.Background()

	opt := &github.RepositoryContentGetOptions{Ref: branch}
	_, directoryContent, resp, err := c.client.Repositories.GetContents(ctx, c.cfg.GitConfig.OrgName, repo, path, opt)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, err
	}
	if directoryContent == nil {
		return nil, nil
	}
	for _, file := range directoryContent {
		files = append(files, file.GetName())
	}
	return files, nil
}

func (c *GithubClientImpl) GetFile(repo string, branch string, path string) (*CommitFile, error) {
	var commitFile CommitFile

	ctx := context.Background()
	opt := &github.RepositoryContentGetOptions{Ref: branch}
	fileContent, _, resp, err := c.client.Repositories.GetContents(ctx, c.cfg.GitConfig.OrgName, repo, path, opt)
	if err != nil {
		return &commitFile, err
	}
	if resp.StatusCode != 200 {
		return &commitFile, err
	}
	if fileContent == nil {
		return &commitFile, nil
	}
	filePath := fileContent.GetPath()
	commitFile.Path = &filePath
	fileContentString, err := fileContent.GetContent()
	if err != nil {
		return &commitFile, err
	}
	commitFile.Content = &fileContentString

	return &commitFile, nil
}

func (c *GithubClientImpl) SetWebhook() error {
	// TODO: validate secret
	ctx := context.Background()
	hook := &github.Hook{
		Config: map[string]interface{}{
			"url":          c.cfg.GitConfig.WebhookURL,
			"content_type": "json",
			"secret":       c.cfg.GitConfig.WebhookSecret, // TODO webhook from k8s secret

		},
		Events: []string{"push", "pull_request"},
		Active: github.Bool(true),
	}
	if c.cfg.GitConfig.OrgLevelWebhook {
		respHook, ok := isOrgWebhookEnabled(ctx, c)
		if !ok {
			retHook, resp, err := c.client.Organizations.CreateHook(ctx, c.cfg.GitConfig.OrgName, hook)
			if err != nil {
				return err
			}
			if resp.StatusCode != 201 {
				return fmt.Errorf("failed to create org level webhhok, API returned %d", resp.StatusCode)
			}
			c.hooks = append(c.hooks, retHook)
		} else {
			c.hooks = append(c.hooks, respHook)
		}

		return nil
	} else {
		for _, repo := range strings.Split(c.cfg.GitConfig.RepoList, ",") {
			respHook, ok := isRepoWebhookEnabled(ctx, c, repo)
			if !ok {
				_, resp, err := c.client.Repositories.CreateHook(ctx, c.cfg.GitConfig.OrgName, repo, hook)
				if err != nil {
					return err
				}

				if resp.StatusCode != 201 {
					return fmt.Errorf("failed to create repo level webhhok for %s, API returned %d", repo, resp.StatusCode)
				}
				c.hooks = append(c.hooks, hook)
			}
			c.hooks = append(c.hooks, respHook)
		}
	}

	return nil
}

func (c *GithubClientImpl) UnsetWebhook() error {
	ctx := context.Background()

	for _, hook := range c.hooks {
		if c.cfg.GitConfig.OrgLevelWebhook {

			resp, err := c.client.Organizations.DeleteHook(ctx, c.cfg.GitConfig.OrgName, *hook.ID)

			if err != nil {
				return err
			}

			if resp.StatusCode != 204 {
				return fmt.Errorf("failed to delete org level webhhok, API call returned %d", resp.StatusCode)
			}

		} else {
			for _, repo := range strings.Split(c.cfg.GitConfig.RepoList, ",") {
				resp, err := c.client.Repositories.DeleteHook(ctx, c.cfg.GitConfig.OrgName, repo, *hook.ID)

				if err != nil {
					return fmt.Errorf("failed to delete repo level webhhok for %s, API call returned %d. %s", repo, resp.StatusCode, err)
				}

				if resp.StatusCode != 204 {
					return fmt.Errorf("failed to delete repo level webhhok for %s, API call returned %d", repo, resp.StatusCode)
				}
			}
		}
	}

	return nil
}

func (c *GithubClientImpl) HandlePayload(request *http.Request, secret []byte) (*WebhookPayload, error) {

	var webhookPayload *WebhookPayload

	payload, err := github.ValidatePayload(request, secret)
	if err != nil {
		return nil, err
	}

	event, err := github.ParseWebHook(github.WebHookType(request), payload)
	if err != nil {
		return nil, err
	}

	switch e := event.(type) {
	case *github.PingEvent:
		webhookPayload = &WebhookPayload{
			Event: "ping",
			Repo:  e.GetRepo().GetFullName(),
		}
	case *github.PushEvent:
		webhookPayload = &WebhookPayload{
			Event:     "push",
			Repo:      e.GetRepo().GetName(),
			Branch:    strings.TrimPrefix(e.GetRef(), "refs/heads/"),
			Commit:    e.GetHeadCommit().GetSHA(),
			User:      e.GetSender().GetName(),
			UserEmail: e.GetSender().GetEmail(),
		}
	case *github.PullRequestEvent:
		webhookPayload = &WebhookPayload{
			Event:            "pull_request",
			Repo:             e.GetRepo().GetName(),
			Branch:           e.GetPullRequest().GetHead().GetRef(),
			Commit:           e.GetPullRequest().GetHead().GetSHA(),
			User:             e.GetSender().GetName(),
			UserEmail:        e.GetSender().GetEmail(),
			PullRequestTitle: e.GetPullRequest().GetTitle(),
			PullRequestURL:   e.GetPullRequest().GetURL(),
			DestBranch:       e.GetPullRequest().GetBase().GetRef(),
		}
	}

	return webhookPayload, nil

}
