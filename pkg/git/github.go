package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/rookout/piper/pkg/conf"

	"github.com/google/go-github/v52/github"
)

type ClientImpl struct {
	client *github.Client
	cfg    *conf.Config
}

func NewGithubClient(cfg *conf.Config) Client {
	ctx := context.Background()

	client := github.NewTokenClient(ctx, cfg.GitConfig.Token)
	ValidatePermissions(ctx, client, cfg)

	return &ClientImpl{
		client: client,
		cfg:    cfg,
	}
}

func ValidatePermissions(ctx context.Context, client *github.Client, cfg *conf.Config) error {

	if cfg.GitConfig.OrgLevelWebhook {
		return nil
	} else {
		for _, repo := range strings.Split(cfg.GitConfig.RepoList, ",") {
			_, _, err := client.Repositories.Get(ctx, cfg.GitConfig.OrgName, repo)
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
}

func (c ClientImpl) ListFiles(repo string, branch string, path string) ([]string, error) {
	var files []string
	ctx := context.Background()

	opt := &github.RepositoryContentGetOptions{Ref: branch}
	_, directoryContent, resp, err := c.client.Repositories.GetContents(ctx, c.cfg.GitConfig.OrgName, repo, path, opt)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, err
	}
	if directoryContent == nil {
		return nil, nil
	}
	for _, file := range directoryContent {
		files = append(files, file.GetPath())
	}
	return files, nil
}

func (c ClientImpl) GetFile(repo string, branch string, path string) (*CommitFile, error) {
	var commitFile CommitFile

	ctx := context.Background()
	opt := &github.RepositoryContentGetOptions{Ref: branch}
	fileContent, _, resp, err := c.client.Repositories.GetContents(ctx, c.cfg.GitConfig.OrgName, repo, path, opt)
	if err != nil {
		return &commitFile, err
	}
	if resp.StatusCode != 200 {
		return &commitFile, err
	}
	if fileContent == nil {
		return &commitFile, nil
	}
	*commitFile.Path = fileContent.GetPath()
	*commitFile.Content, err = fileContent.GetContent()
	if err != nil {
		return &commitFile, err
	}

	return &commitFile, nil
}

func (c ClientImpl) SetWebhook() error {
	ctx := context.Background()
	hook := &github.Hook{
		Config: map[string]interface{}{
			"url":          "@123",
			"content_type": "json",
			//"secret":       "123", // TODO webhook secret

		},
		Events: []string{"push", "pull_request"},
		Active: github.Bool(true),
	}
	if c.cfg.GitConfig.OrgLevelWebhook {
		if !isOrgWebhookEnabled(ctx, &c) {
			_, resp, err := c.client.Organizations.CreateHook(ctx, c.cfg.GitConfig.OrgName, hook)
			if err != nil {
				return err
			}
			if resp.StatusCode != 200 {
				return fmt.Errorf("failed to create org level webhhok, API returned %d", resp.StatusCode)
			}
		}
		return nil
	} else {
		for _, repo := range strings.Split(c.cfg.GitConfig.RepoList, ",") {
			if !isRepoWebhookEnabled(ctx, &c, repo) {
				_, resp, err := c.client.Repositories.CreateHook(ctx, c.cfg.GitConfig.OrgName, repo, hook)
				if err != nil {
					return err
				}

				if resp.StatusCode != 200 {
					return fmt.Errorf("failed to create repo level webhhok for %s, API returned %d", repo, resp.StatusCode)
				}
			}
		}
	}

	return nil
}

func (c ClientImpl) UnsetWebhook() error {
	//TODO implement me
	panic("implement me")
}
