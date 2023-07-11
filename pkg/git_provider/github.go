package git_provider

import (
	"context"
	"fmt"
	"github.com/rookout/piper/pkg/utils"
	"log"
	"net/http"
	"strings"

	"github.com/rookout/piper/pkg/conf"

	"github.com/google/go-github/v52/github"
)

type GithubClientImpl struct {
	client *github.Client
	cfg    *conf.GlobalConfig
	hooks  []*github.Hook
}

func NewGithubClient(cfg *conf.GlobalConfig) (Client, error) {
	ctx := context.Background()

	client := github.NewTokenClient(ctx, cfg.GitProviderConfig.Token)
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

func (c *GithubClientImpl) ListFiles(ctx *context.Context, repo string, branch string, path string) ([]string, error) {
	var files []string

	opt := &github.RepositoryContentGetOptions{Ref: branch}
	_, directoryContent, resp, err := c.client.Repositories.GetContents(*ctx, c.cfg.GitProviderConfig.OrgName, repo, path, opt)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("github provider returned %d: failed to get contents of %s/%s%s", resp.StatusCode, repo, branch, path)
	}
	if directoryContent == nil {
		return nil, nil
	}
	for _, file := range directoryContent {
		files = append(files, file.GetName())
	}
	return files, nil
}

func (c *GithubClientImpl) GetFile(ctx *context.Context, repo string, branch string, path string) (*CommitFile, error) {
	var commitFile CommitFile

	opt := &github.RepositoryContentGetOptions{Ref: branch}
	fileContent, _, resp, err := c.client.Repositories.GetContents(*ctx, c.cfg.GitProviderConfig.OrgName, repo, path, opt)
	if err != nil {
		return &commitFile, err
	}
	if resp.StatusCode == 404 {
		log.Printf("File %s not found in repo %s branch %s", path, repo, branch)
		return nil, nil
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

func (c *GithubClientImpl) GetFiles(ctx *context.Context, repo string, branch string, paths []string) ([]*CommitFile, error) {
	var commitFiles []*CommitFile
	for _, path := range paths {
		file, err := c.GetFile(ctx, repo, branch, path)
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

func (c *GithubClientImpl) SetWebhook() error {
	ctx := context.Background()
	hook := &github.Hook{
		Config: map[string]interface{}{
			"url":          c.cfg.GitProviderConfig.WebhookURL,
			"content_type": "json",
			"secret":       c.cfg.GitProviderConfig.WebhookSecret,
		},
		Events: []string{"push", "pull_request", "create"},
		Active: github.Bool(true),
	}
	if c.cfg.GitProviderConfig.OrgLevelWebhook {
		respHook, ok := isOrgWebhookEnabled(ctx, c)
		if !ok {
			createdHook, resp, err := c.client.Organizations.CreateHook(
				ctx,
				c.cfg.GitProviderConfig.OrgName,
				hook,
			)
			if err != nil {
				return err
			}
			if resp.StatusCode != 201 {
				return fmt.Errorf("failed to create org level webhhok, API returned %d", resp.StatusCode)
			}
			log.Printf("edited webhook of type %s for %s :%s\n", createdHook.GetType(), c.cfg.GitProviderConfig.OrgName, createdHook.GetURL())
			c.hooks = append(c.hooks, createdHook)
		} else {
			updatedHook, resp, err := c.client.Organizations.EditHook(
				ctx,
				c.cfg.GitProviderConfig.OrgName,
				respHook.GetID(),
				hook,
			)
			if err != nil {
				return err
			}
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf(
					"failed to update org level webhhok for %s, API returned %d",
					c.cfg.GitProviderConfig.OrgName,
					resp.StatusCode,
				)
			}
			log.Printf("edited webhook of type %s for %s :%s\n", updatedHook.GetType(), c.cfg.GitProviderConfig.OrgName, updatedHook.GetURL())
			c.hooks = append(c.hooks, updatedHook)
		}

		return nil
	} else {
		for _, repo := range strings.Split(c.cfg.GitProviderConfig.RepoList, ",") {
			respHook, ok := isRepoWebhookEnabled(ctx, c, repo)
			if !ok {
				createdHook, resp, err := c.client.Repositories.CreateHook(ctx, c.cfg.GitProviderConfig.OrgName, repo, hook)
				if err != nil {
					return err
				}

				if resp.StatusCode != 201 {
					return fmt.Errorf("failed to create repo level webhhok for %s, API returned %d", repo, resp.StatusCode)
				}
				log.Printf("created webhook of type %s for %s :%s\n", createdHook.GetType(), repo, createdHook.GetURL())
				c.hooks = append(c.hooks, createdHook)
			} else {
				updatedHook, resp, err := c.client.Repositories.EditHook(ctx, c.cfg.GitProviderConfig.OrgName, repo, respHook.GetID(), hook)
				if err != nil {
					return err
				}
				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("failed to update repo level webhhok for %s, API returned %d", repo, resp.StatusCode)
				}
				log.Printf("edited webhook of type %s for %s :%s\n", updatedHook.GetType(), repo, updatedHook.GetURL())
				c.hooks = append(c.hooks, updatedHook)
			}
		}
	}

	return nil
}

func (c *GithubClientImpl) UnsetWebhook(ctx *context.Context) error {

	for _, hook := range c.hooks {
		if c.cfg.GitProviderConfig.OrgLevelWebhook {

			resp, err := c.client.Organizations.DeleteHook(*ctx, c.cfg.GitProviderConfig.OrgName, *hook.ID)

			if err != nil {
				return err
			}

			if resp.StatusCode != 204 {
				return fmt.Errorf("failed to delete org level webhhok, API call returned %d", resp.StatusCode)
			}

		} else {
			for _, repo := range strings.Split(c.cfg.GitProviderConfig.RepoList, ",") {
				resp, err := c.client.Repositories.DeleteHook(*ctx, c.cfg.GitProviderConfig.OrgName, repo, *hook.ID)

				if err != nil {
					return fmt.Errorf("failed to delete repo level webhhok for %s, API call returned %d. %s", repo, resp.StatusCode, err)
				}

				if resp.StatusCode != 204 {
					return fmt.Errorf("failed to delete repo level webhhok for %s, API call returned %d", repo, resp.StatusCode)
				}
			}
		}
		log.Printf("removed hook:%s\n", hook.GetURL()) // INFO
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
			Action:    e.GetAction(),
			Repo:      e.GetRepo().GetName(),
			Branch:    strings.TrimPrefix(e.GetRef(), "refs/heads/"),
			Commit:    e.GetHeadCommit().GetID(),
			User:      e.GetSender().GetLogin(),
			UserEmail: e.GetHeadCommit().GetAuthor().GetEmail(),
		}
	case *github.PullRequestEvent:
		webhookPayload = &WebhookPayload{
			Event:            "pull_request",
			Action:           e.GetAction(),
			Repo:             e.GetRepo().GetName(),
			Branch:           e.GetPullRequest().GetHead().GetRef(),
			Commit:           e.GetPullRequest().GetHead().GetSHA(),
			User:             e.GetPullRequest().GetUser().GetLogin(),
			UserEmail:        e.GetSender().GetEmail(), // e.GetPullRequest().GetUser().GetEmail() Not working. GitHub missing email for PR events in payload.
			PullRequestTitle: e.GetPullRequest().GetTitle(),
			PullRequestURL:   e.GetPullRequest().GetHTMLURL(),
			DestBranch:       e.GetPullRequest().GetBase().GetRef(),
			Labels:           e.GetPullRequest().Labels,
		}
	case *github.CreateEvent:
		webhookPayload = &WebhookPayload{
			Event:     "create",
			Action:    e.GetRefType(), // Possible values are: "repository", "branch", "tag".
			Repo:      e.GetRepo().GetName(),
			Branch:    e.GetRef(),
			Commit:    e.GetRef(),
			User:      e.GetSender().GetLogin(),
			UserEmail: e.GetSender().GetEmail(),
		}
	}

	return webhookPayload, nil

}

func (c *GithubClientImpl) SetStatus(ctx *context.Context, repo *string, commit *string, linkURL *string, status *string, message *string) error {
	if !utils.ValidateHTTPFormat(*linkURL) {
		return fmt.Errorf("invalid linkURL")
	}
	repoStatus := &github.RepoStatus{
		State:       status, // pending, success, error, or failure.
		TargetURL:   linkURL,
		Description: utils.SPtr(fmt.Sprintf("Workflow %s %s", *status, *message)),
		Context:     utils.SPtr("Piper/ArgoWorkflows"),
		AvatarURL:   utils.SPtr("https://argoproj.github.io/argo-workflows/assets/logo.png"),
	}
	_, resp, err := c.client.Repositories.CreateStatus(*ctx, c.cfg.OrgName, *repo, *commit, repoStatus)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to set status on repo:%s, commit:%s, API call returned %d", *repo, *commit, resp.StatusCode)
	}

	log.Printf("successfully set status on repo:%s commit: %s to status: %s\n", *repo, *commit, *status)
	return nil
}
