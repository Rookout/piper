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
}

func NewGithubClient(cfg *conf.GlobalConfig) (Client, error) {
	ctx := context.Background()

	client := github.NewTokenClient(ctx, cfg.GitProviderConfig.Token)
	err := ValidatePermissions(ctx, client, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to validate permissions: %v", err)
	}

	user, resp, err := client.Users.Get(context.Background(), cfg.OrgName)
	if err != nil {
		return nil, fmt.Errorf("failed to get org id: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get organization id %s", resp.Status)
	}

	cfg.OrgID = user.GetID()

	fmt.Printf("Org ID is: %d\n", cfg.OrgID)

	return &GithubClientImpl{
		client: client,
		cfg:    cfg,
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

func (c *GithubClientImpl) SetWebhook(ctx *context.Context, repo *string) (*HookWithStatus, error) {
	if c.cfg.OrgLevelWebhook && repo != nil {
		return nil, fmt.Errorf("trying to set repo scope. repo: %s", *repo)
	}

	hookConf := &github.Hook{
		Config: map[string]interface{}{
			"url":          c.cfg.GitProviderConfig.WebhookURL,
			"content_type": "json",
			"secret":       c.cfg.GitProviderConfig.WebhookSecret,
		},
		Events: []string{"push", "pull_request", "create", "release"},
		Active: github.Bool(true),
	}

	if repo == nil {
		respHook, ok := isOrgWebhookEnabled(*ctx, c)
		if !ok {
			createdHook, resp, err := c.client.Organizations.CreateHook(
				*ctx,
				c.cfg.GitProviderConfig.OrgName,
				hookConf,
			)
			if err != nil {
				return nil, err
			}
			if resp.StatusCode != 201 {
				return nil, fmt.Errorf("failed to create org level webhhok, API returned %d", resp.StatusCode)
			}
			log.Printf("edited webhook of type %s for %s name: %s\n", createdHook.GetType(), c.cfg.GitProviderConfig.OrgName, createdHook.Config["url"])
			hookID := createdHook.GetID()
			return &HookWithStatus{HookID: hookID, HealthStatus: true, RepoName: repo}, nil
		} else {
			updatedHook, resp, err := c.client.Organizations.EditHook(
				*ctx,
				c.cfg.GitProviderConfig.OrgName,
				respHook.GetID(),
				hookConf,
			)
			if err != nil {
				return nil, err
			}
			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf(
					"failed to update org level webhhok for %s, API returned %d",
					c.cfg.GitProviderConfig.OrgName,
					resp.StatusCode,
				)
			}
			log.Printf("edited webhook of type %s for %s: %s\n", updatedHook.GetType(), c.cfg.GitProviderConfig.OrgName, updatedHook.Config["url"])
			hookID := updatedHook.GetID()
			return &HookWithStatus{HookID: hookID, HealthStatus: true, RepoName: repo}, nil
		}
	} else {
		respHook, ok := isRepoWebhookEnabled(*ctx, c, *repo)
		if !ok {
			createdHook, resp, err := c.client.Repositories.CreateHook(*ctx, c.cfg.GitProviderConfig.OrgName, *repo, hookConf)
			if err != nil {
				return nil, err
			}

			if resp.StatusCode != 201 {
				return nil, fmt.Errorf("failed to create repo level webhhok for %s, API returned %d", *repo, resp.StatusCode)
			}
			log.Printf("created webhook of type %s for %s: %s\n", createdHook.GetType(), *repo, createdHook.Config["url"])
			hookID := createdHook.GetID()
			return &HookWithStatus{HookID: hookID, HealthStatus: true, RepoName: repo}, nil
		} else {
			updatedHook, resp, err := c.client.Repositories.EditHook(*ctx, c.cfg.GitProviderConfig.OrgName, *repo, respHook.GetID(), hookConf)
			if err != nil {
				return nil, err
			}
			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("failed to update repo level webhhok for %s, API returned %d", *repo, resp.StatusCode)
			}
			log.Printf("edited webhook of type %s for %s: %s\n", updatedHook.GetType(), *repo, updatedHook.Config["url"])
			hookID := updatedHook.GetID()
			return &HookWithStatus{HookID: hookID, HealthStatus: true, RepoName: repo}, nil
		}

	}
}

func (c *GithubClientImpl) UnsetWebhook(ctx *context.Context, hook *HookWithStatus) error {

	if hook.RepoName == nil {

		resp, err := c.client.Organizations.DeleteHook(*ctx, c.cfg.GitProviderConfig.OrgName, hook.HookID)

		if err != nil {
			return err
		}

		if resp.StatusCode != 204 {
			return fmt.Errorf("failed to delete org level webhhok, API call returned %d", resp.StatusCode)
		}
		log.Printf("removed org webhook, hookID :%d\n", hook.HookID) // INFO
	} else {
		resp, err := c.client.Repositories.DeleteHook(*ctx, c.cfg.GitProviderConfig.OrgName, *hook.RepoName, hook.HookID)

		if err != nil {
			return fmt.Errorf("failed to delete repo level webhhok for %s, API call returned %d. %s", *hook.RepoName, resp.StatusCode, err)
		}

		if resp.StatusCode != 204 {
			return fmt.Errorf("failed to delete repo level webhhok for %s, API call returned %d", *hook.RepoName, resp.StatusCode)
		}
		log.Printf("removed repo webhook, repo:%s hookID :%d\n", *hook.RepoName, hook.HookID) // INFO
	}

	return nil
}

func (c *GithubClientImpl) HandlePayload(ctx *context.Context, request *http.Request, secret []byte) (*WebhookPayload, error) {
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
			Event:   "ping",
			Repo:    e.GetRepo().GetFullName(),
			HookID:  e.GetHookID(),
			OwnerID: e.GetSender().GetID(),
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
			OwnerID:   e.GetSender().GetID(),
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
			Labels:           c.extractLabelNames(e.GetPullRequest().Labels),
			OwnerID:          e.GetSender().GetID(),
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
			OwnerID:   e.GetSender().GetID(),
		}
	case *github.ReleaseEvent:
		commitSHA, _err := c.refToSHA(ctx, e.GetRelease().GetName(), e.GetRepo().GetName())
		if _err != nil {
			return webhookPayload, _err
		}
		webhookPayload = &WebhookPayload{
			Event:     "release",
			Action:    e.GetAction(), // "created", "edited", "deleted", or "prereleased".
			Repo:      e.GetRepo().GetName(),
			Branch:    e.GetRelease().GetTagName(),
			Commit:    *commitSHA,
			User:      e.GetSender().GetLogin(),
			UserEmail: e.GetSender().GetEmail(),
			OwnerID:   e.GetSender().GetID(),
		}
	}

	if c.cfg.EnforceOrgBelonging && webhookPayload.OwnerID != c.cfg.OrgID {
		return nil, fmt.Errorf("webhook send from non organizational member")
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

func (c *GithubClientImpl) PingHook(ctx *context.Context, hook *HookWithStatus) error {
	if c.cfg.OrgLevelWebhook && hook.RepoName != nil {
		return fmt.Errorf("trying to ping repo scope webhook while configured for org level webhook. repo: %s", *hook.RepoName)
	}
	if hook.RepoName == nil {
		resp, err := c.client.Organizations.PingHook(*ctx, c.cfg.OrgName, hook.HookID)
		if err != nil {
			return err
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("unable to find organization webhook for hookID: %d", hook.HookID)
		}
	} else {
		resp, err := c.client.Repositories.PingHook(*ctx, c.cfg.GitProviderConfig.OrgName, *hook.RepoName, hook.HookID)
		if err != nil {
			return err
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("unable to find repo webhook for repo:%s hookID: %d", *hook.RepoName, hook.HookID)
		}
	}

	return nil
}

func (c *GithubClientImpl) refToSHA(ctx *context.Context, ref string, repo string) (*string, error) {
	respSHA, resp, err := c.client.Repositories.GetCommitSHA1(*ctx, c.cfg.OrgName, repo, ref, "")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to set status on repo:%s, commit:%s, API call returned %d", repo, ref, resp.StatusCode)
	}

	log.Printf("resolved ref: %s to SHA: %s", ref, respSHA)
	return &respSHA, nil
}

func (c *GithubClientImpl) extractLabelNames(labels []*github.Label) []string {
	var returnLabelsList []string
	for _, label := range labels {
		returnLabelsList = append(returnLabelsList, *label.Name)
	}
	return returnLabelsList
}
