package git

import (
	"context"

	"github.com/rookout/piper/pkg/conf"

	github "github.com/google/go-github/v52/github"
)

type Client struct {
	client *github.Client
	cfg    *conf.Config
}

func (c Client) New() Client {
	ctx := context.Background()

	client := github.NewTokenClient(ctx, c.cfg.GitConfig.Token)

	return Client{
		client: client,
		cfg:    c.cfg,
	}
}

func (c Client) ListFiles(repo string, branch string, path string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (c Client) GetFile(repo string, branch string, path string) (CommitFile, error) {
	//TODO implement me
	panic("implement me")
}

func (c Client) SetWebhook(org string, repo string) error {
	//TODO implement me
	panic("implement me")
}

func (c Client) UnsetWebhook(org string, repo string) error {
	//TODO implement me
	panic("implement me")
}
