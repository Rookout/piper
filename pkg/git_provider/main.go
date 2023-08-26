package git_provider

import (
	"fmt"
	"github.com/rookout/piper/pkg/conf"
)

func NewGitProviderClient(cfg *conf.GlobalConfig) (Client, error) {

	switch cfg.GitProviderConfig.Provider {
	case "github":
		gitClient, err := NewGithubClient(cfg)
		if err != nil {
			return nil, err
		}
		return gitClient, nil
	case "bitbucket":
		gitClient, err := NewBitbucketServerClient(cfg)
		if err != nil {
			return nil, err
		}
		return gitClient, nil
	}

	return nil, fmt.Errorf("didn't find matching git provider %s", cfg.GitProviderConfig.Provider)
}
