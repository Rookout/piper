package git

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/v52/github"
	"github.com/rookout/piper/pkg/conf"
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

func GetScopes(ctx context.Context, client *github.Client) ([]string, error) {
	// Make a request to the "Get the authenticated user" endpoint
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	resp, err := client.Do(context.Background(), req, nil)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check the "X-OAuth-Scopes" header to get the token scopes
	scopes := resp.Header.Get("X-OAuth-Scopes")
	fmt.Println("Github Token Scopes are:", scopes)

	scopes = strings.ReplaceAll(scopes, " ", "")
	return strings.Split(scopes, ","), nil

}

func ValidateListAInListB(listA, listB []string) bool {
	for _, element := range listA {
		found := false
		for _, b := range listB {
			if element == b {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func ValidatePermissions(ctx context.Context, client *github.Client, cfg *conf.Config) error {

	orgScopes := []string{"admin:org_hook"}
	repoAdminScopes := []string{"admin:repo_hook"}
	repoGranularScopes := []string{"write:repo_hook", "read:repo_hook"}

	scopes, err := GetScopes(ctx, client)

	if err != nil {
		return fmt.Errorf("failed to get scopes: %v", err)
	}
	if len(scopes) == 0 {
		return fmt.Errorf("permissions error: no scopes found for the github client")
	}

	if cfg.GitConfig.OrgLevelWebhook {
		if ValidateListAInListB(orgScopes, scopes) {
			return nil
		}
		return fmt.Errorf("permissions error: %v is not a valid scope for the org level permissions", scopes)
	}

	if ValidateListAInListB(repoAdminScopes, scopes) {
		return nil
	}
	if ValidateListAInListB(repoGranularScopes, scopes) {
		return nil
	}

	return fmt.Errorf("permissions error: %v is not a valid scope for the repo level permissions", scopes)
}
