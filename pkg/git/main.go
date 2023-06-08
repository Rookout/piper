package git

import (
	"github.com/rookout/piper/pkg/conf"
)

func NewGitProviderClient(cfg *conf.Config) Client {

	switch cfg.GitConfig.Provider {
	case "github":
		gitClient, err := NewGithubClient(cfg)
		if err != nil {
			panic(err)
		}
		return gitClient
	}

	return nil
}
