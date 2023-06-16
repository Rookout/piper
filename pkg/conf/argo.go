package conf

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type ArgoConfig struct {
	ArgoToken   string `envconfig:"ARGO_WORKFLOWS_TOKEN" required:"true"`
	ArgoAddress string `envconfig:"ARGO_WORKFLOWS_ADDRESS" required:"true"`
	CreateCRD   bool   `envconfig:"ARGO_WORKFLOWS_CREATE_CRD" default:"true"`
	Namespace   string `envconfig:"ARGO_WORKFLOWS_NAMESPACE" default:"default"`
	KubeConfig  string `envconfig:"KUBE_CONFIG" default:""`
}

func (cfg *ArgoConfig) ArgoConfLoad() error {
	err := envconfig.Process("", cfg)
	if err != nil {
		return fmt.Errorf("failed to load the Argo configuration, error: %v", err)
	}

	return nil
}
