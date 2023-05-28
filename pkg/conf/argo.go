package conf

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type ArgoConfig struct {
	ArgoToken   string `envconfig:"ARGO_WORKFLOWS_TOKEN" required:"true"`
	ArgoAddress string `envconfig:"ARGO_WORKFLOWS_ADDRESS" required:"true"`
	CreateCRD   bool   `envconfig:"GIT_WEBHOOK_REPO_LIST" default:"true"`
}

func (cfg *GitConfig) ArgoConfLoad() error {
	err := envconfig.Process("", cfg)
	if err != nil {
		return fmt.Errorf("failed to load the Git provider configuration, error: %v", err)
	}

	return nil
}
