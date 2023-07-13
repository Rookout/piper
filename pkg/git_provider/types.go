package git_provider

import (
	"context"
	"github.com/google/go-github/v52/github"
	"net/http"
)

type HookWithStatus struct {
	HookID       *int64
	HealthStatus bool
	RepoName     *string
}

type CommitFile struct {
	Path    *string `json:"path"`
	Content *string `json:"content"`
}

type WebhookPayload struct {
	Event            string          `json:"event"`
	Action           string          `json:"action"`
	Repo             string          `json:"repoName"`
	Branch           string          `json:"branch"`
	Commit           string          `json:"commit"`
	User             string          `json:"user"`
	UserEmail        string          `json:"user_email"`
	PullRequestURL   string          `json:"pull_request_url"`
	PullRequestTitle string          `json:"pull_request_title"`
	DestBranch       string          `json:"dest_branch"`
	Labels           []*github.Label `json:"labels"`
	HookID           int64           `json:"hookID"`
}

type Client interface {
	ListFiles(ctx *context.Context, repo string, branch string, path string) ([]string, error)
	GetFile(ctx *context.Context, repo string, branch string, path string) (*CommitFile, error)
	GetFiles(ctx *context.Context, repo string, branch string, paths []string) ([]*CommitFile, error)
	SetWebhook(ctx *context.Context, repo *string) (*HookWithStatus, error)
	UnsetWebhook(ctx *context.Context, hook *HookWithStatus) error
	HandlePayload(request *http.Request, secret []byte) (*WebhookPayload, error)
	SetStatus(ctx *context.Context, repo *string, commit *string, linkURL *string, status *string, message *string) error
	PingHook(ctx *context.Context, hook *HookWithStatus) error
}
