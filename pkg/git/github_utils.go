package git

import (
	"context"

	"github.com/google/go-github/v52/github"
)

func isOrgWebhookEnabled(ctx context.Context, c *ClientImpl) bool {
	hooks, resp, err := c.client.Organizations.ListHooks(ctx, c.cfg.GitConfig.OrgName, &github.ListOptions{})
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	if hooks == nil || len(hooks) == 0 {
		return false
	}
	for _, hook := range hooks {
		if hook.GetActive() && hook.GetName() == "piper" {
			return true
		}
	}
	return false
}

func isRepoWebhookEnabled(ctx context.Context, c *ClientImpl, repo string) bool {
	hooks, resp, err := c.client.Repositories.ListHooks(ctx, c.cfg.GitConfig.OrgName, repo, &github.ListOptions{})
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	if hooks == nil || len(hooks) == 0 {
		return false
	}

	for _, hook := range hooks {
		if hook.GetActive() && hook.GetName() == "piper" {
			return true
		}
	}

	return false
}
