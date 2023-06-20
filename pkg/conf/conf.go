package conf

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type GlobalConfig struct {
	GitProviderConfig
	WorkflowServerConfig
	RookoutConfig
	WorkflowsConfig
}

func (cfg *GlobalConfig) Load() error {
	err := envconfig.Process("", cfg)
	if err != nil {
		return fmt.Errorf("failed to load the configuration, error: %v", err)
	}

	return nil
}

func LoadConfig() (*GlobalConfig, error) {
	cfg := new(GlobalConfig)

	err := cfg.Load()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
