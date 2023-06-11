package git

import (
	"context"

	"github.com/google/go-github/v52/github"
)

func isOrgWebhookEnabled(ctx context.Context, c *GithubClientImpl) (*github.Hook, bool) {
	emptyHook := github.Hook{}
	hooks, resp, err := c.client.Organizations.ListHooks(ctx, c.cfg.GitConfig.OrgName, &github.ListOptions{})
	if err != nil {
		return &emptyHook, false
	}
	if resp.StatusCode != 200 {
		return &emptyHook, false
	}
	if hooks == nil || len(hooks) == 0 {
		return &emptyHook, false
	}
	for _, hook := range hooks {
		if hook.GetActive() && hook.GetName() == "web" && hook.Config["url"] == c.cfg.GitConfig.WebhookURL {
			return hook, true
		}
	}
	return &emptyHook, false
}

func isRepoWebhookEnabled(ctx context.Context, c *GithubClientImpl, repo string) (*github.Hook, bool) {
	emptyHook := github.Hook{}
	hooks, resp, err := c.client.Repositories.ListHooks(ctx, c.cfg.GitConfig.OrgName, repo, &github.ListOptions{})
	if err != nil {
		return &emptyHook, false
	}
	if resp.StatusCode != 200 {
		return &emptyHook, false
	}
	if hooks == nil || len(hooks) == 0 {
		return &emptyHook, false
	}

	for _, hook := range hooks {
		if hook.GetActive() && hook.GetName() == "web" && hook.Config["url"] == c.cfg.GitConfig.WebhookURL {
			return hook, true
		}
	}

	return &emptyHook, false
}
