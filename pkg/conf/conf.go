package conf

import (
	"fmt"

	"github.com/rookout/piper/pkg/git"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	GitConfig
}

type Clients struct {
	Git git.Client
}

func (cfg *Config) Load() error {
	err := envconfig.Process("", cfg)
	if err != nil {
		return fmt.Errorf("failed to load the configuration, error: %v", err)
	}

	return nil
}

func LoadConfig() (*Config, error) {
	cfg := new(Config)

	err := cfg.Load()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
