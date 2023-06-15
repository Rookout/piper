package conf

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type RookoutConfig struct {
	Token        string `envconfig:"ROOKOUT_TOKEN" default:"true"`
	Labels       string `envconfig:"ROOKOUT_LABELS" default:"service:piper"`
	RemoteOrigin string `envconfig:"ROOKOUT_REMOTE_ORIGIN" default:"https://github.com/Rookout/piper.git"`
}

func (cfg *GitConfig) RookoutConfLoad() error {
	err := envconfig.Process("", cfg)
	if err != nil {
		return fmt.Errorf("failed to load the Git provider configuration, error: %v", err)
	}

	return nil
}
