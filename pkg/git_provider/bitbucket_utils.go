package git_provider

import (
	"fmt"
	bitbucket "github.com/gfleury/go-bitbucket-v1"
	"github.com/rookout/piper/pkg/conf"
	"golang.org/x/net/context"
)

func BitbucketValidatePermissions(ctx context.Context, client *bitbucket.APIClient, cfg *conf.GlobalConfig) error {
	resp, err := client.DefaultApi.GetRepositories(cfg.GitProviderConfig.OrgName)

	if resp.StatusCode != 200 {
		return fmt.Errorf("bitbucket get repos returned status %s", resp.Status)
	}

	return err
}
