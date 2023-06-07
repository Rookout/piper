package git

type CommitFile struct {
	Path    *string `json:"path"`
	Content *string `json:"content"`
}

type Client interface {
	ListFiles(repo string, branch string, path string) ([]string, error)
	GetFile(repo string, branch string, path string) (*CommitFile, error)
	SetWebhook() error
	UnsetWebhook() error
}
