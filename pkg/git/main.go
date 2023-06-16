package git

import (
	"fmt"
	"github.com/rookout/piper/pkg/conf"
)

func NewGitProviderClient(cfg *conf.Config) (Client, error) {

	switch cfg.GitConfig.Provider {
	case "github":
		gitClient, err := NewGithubClient(cfg)
		if err != nil {
			return nil, err
		}
		return gitClient, nil
	}

	return nil, fmt.Errorf("didn't find matching git provider %s", cfg.GitConfig.Provider)
}
