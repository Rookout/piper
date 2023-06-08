package git

type CommitFile struct {
	Path    *string `json:"path"`
	Content *string `json:"content"`
}

type WebhookPayload struct {
	Repo             string `json:"repoName"`
	Branch           string `json:"branch"`
	Commit           string `json:"commit"`
	User             string `json:"user"`
	PullRequestUrl   string `json:"pull_request_url"`
	PullRequestTitle string `json:"pull_request_title"`
	DestBranch       string `json:"dest_branch"`
}

type Client interface {
	ListFiles(repo string, branch string, path string) ([]string, error)
	GetFile(repo string, branch string, path string) (*CommitFile, error)
	SetWebhook() error
	UnsetWebhook() error
	ParseWebhookPayload(payload string) (*WebhookPayload, error)
}
