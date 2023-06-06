package git

import (
	"github.com/rookout/piper/pkg/conf"
)

type CommitFile struct {
	Path    *string `json:"path"`
	Content *string `json:"content"`
}

type ClientImpl interface {
	New(cfg *conf.Config) (ClientImpl, error)
	ListFiles(repo string, branch string, path string) ([]string, error)
	GetFile(repo string, branch string, path string) (CommitFile, error)
	SetWebhook(org string, repo string) error
	UnsetWebhook(org string, repo string) error
}
