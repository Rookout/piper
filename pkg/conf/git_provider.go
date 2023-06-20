package conf

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type GitProviderConfig struct {
	Provider        string `envconfig:"GIT_PROVIDER" required:"true"`
	Token           string `envconfig:"GIT_TOKEN" required:"true"`
	OrgName         string `envconfig:"GIT_ORG_NAME" required:"true"`
	OrgLevelWebhook bool   `envconfig:"GIT_ORG_LEVEL_WEBHOOK" default:"false" required:"false"`
	RepoList        string `envconfig:"GIT_WEBHOOK_REPO_LIST" required:"false"`
	WebhookURL      string `envconfig:"GIT_WEBHOOK_URL" required:"false"`
	WebhookSecret   string `envconfig:"GIT_WEBHOOK_SECRET" required:"false"`
}

func (cfg *GitProviderConfig) GitConfLoad() error {
	err := envconfig.Process("", cfg)
	if err != nil {
		return fmt.Errorf("failed to load the Git provider configuration, error: %v", err)
	}

	return nil
}
