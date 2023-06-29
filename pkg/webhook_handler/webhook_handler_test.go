package webhook_handler

import (
	"context"
	"fmt"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/common"
	"github.com/rookout/piper/pkg/git_provider"
	"github.com/rookout/piper/pkg/utils"
	assertion "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// MockGitProvider is a mock implementation of the git_provider.Client interface.
type MockGitProvider struct{}

func GetContent(filename string) *string {
	switch filename {
	case ".workflows/main.yaml":
		return utils.SPtr(`main.yaml`)
	case ".workflows/exit.yaml":
		return utils.SPtr(`exit.yaml`)
	}
	return nil
}

func (m *MockGitProvider) GetFile(ctx *context.Context, repo string, branch string, path string) (*git_provider.CommitFile, error) {
	switch repo {
	case "repo1":
		switch branch {
		case "branch1":
			switch path {
			case ".workflows/main.yaml":
				return &git_provider.CommitFile{
					Path:    &path,
					Content: GetContent(path),
				}, nil
			case ".workflows/exit.yaml":
				return &git_provider.CommitFile{
					Path:    &path,
					Content: GetContent(path),
				}, nil
			}
		}
	}

	return &git_provider.CommitFile{}, nil
}

func (m *MockGitProvider) GetFiles(ctx *context.Context, repo string, branch string, paths []string) ([]*git_provider.CommitFile, error) {
	var commitFiles []*git_provider.CommitFile

	switch repo {
	case "repo1":
		switch branch {
		case "branch1":
			for _, path := range paths {
				toAppend := &git_provider.CommitFile{}
				switch path {
				case ".workflows/main.yaml":
					toAppend = &git_provider.CommitFile{
						Path:    &path,
						Content: GetContent(path),
					}
				case ".workflows/exit.yaml":
					{
						toAppend = &git_provider.CommitFile{
							Path:    &path,
							Content: GetContent(path),
						}
					}
				}

				commitFiles = append(commitFiles, toAppend)
			}
			return commitFiles, nil
		}
	}

	return commitFiles, fmt.Errorf("not found")
}

func (m *MockGitProvider) ListFiles(ctx *context.Context, repo string, branch string, path string) ([]string, error) {
	return nil, nil
}

func (m *MockGitProvider) SetWebhook() error {
	return nil
}

func (m *MockGitProvider) UnsetWebhook() error {
	return nil
}

func (m *MockGitProvider) HandlePayload(request *http.Request, secret []byte) (*git_provider.WebhookPayload, error) {
	return nil, nil
}

func TestPrepareBatchForMatchingTriggers(t *testing.T) {
	assert := assertion.New(t)
	ctx := context.Background()
	tests := []struct {
		name                  string
		triggers              *[]Trigger
		payload               *git_provider.WebhookPayload
		expectedWorkflowBatch []*common.WorkflowsBatch
	}{
		{name: "test1",
			triggers: &[]Trigger{{
				Events:    &[]string{"event1", "event2.action2"},
				Branches:  &[]string{"branch1", "branch2"},
				Templates: &[]string{""},
				OnStart:   &[]string{"main.yaml"},
				OnExit:    &[]string{"exit.yaml"},
				Config:    "default",
			}},
			payload: &git_provider.WebhookPayload{
				Event:            "event1",
				Action:           "",
				Repo:             "repo1",
				Branch:           "branch1",
				Commit:           "commitHSA",
				User:             "piper",
				UserEmail:        "piper@rookout.com",
				PullRequestURL:   "",
				PullRequestTitle: "",
				DestBranch:       "",
			},
			expectedWorkflowBatch: []*common.WorkflowsBatch{
				&common.WorkflowsBatch{
					OnStart: []*git_provider.CommitFile{
						{
							Path:    utils.SPtr(".workflows/main.yaml"),
							Content: GetContent(".workflows/main.yaml"),
						},
					},
					OnExit: []*git_provider.CommitFile{
						{
							Path:    utils.SPtr(".workflows/exit.yaml"),
							Content: GetContent(".workflows/exit.yaml"),
						},
					},
					Templates: []*git_provider.CommitFile{
						&git_provider.CommitFile{
							Path:    nil,
							Content: nil,
						},
					},
					Parameters: &git_provider.CommitFile{
						Path:    nil,
						Content: nil,
					},
					Config:  utils.SPtr("default"),
					Payload: &git_provider.WebhookPayload{},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wh := &WebhookHandlerImpl{
				Triggers: test.triggers,
				Payload:  test.payload,
				clients: &clients.Clients{
					GitProvider: &MockGitProvider{},
				},
			}
			WorkflowsBatches, err := wh.PrepareBatchForMatchingTriggers(&ctx)
			assert.Nil(err)
			for iwf, wf := range WorkflowsBatches {
				for i, _ := range wf.OnStart {
					assert.Equal(*WorkflowsBatches[iwf].OnStart[i].Path, *test.expectedWorkflowBatch[iwf].OnStart[i].Path)
					assert.Equal(*WorkflowsBatches[iwf].OnStart[i].Content, *test.expectedWorkflowBatch[iwf].OnStart[i].Content)
				}
				for j, _ := range wf.OnExit {
					assert.Equal(*WorkflowsBatches[iwf].OnExit[j].Path, *test.expectedWorkflowBatch[iwf].OnExit[j].Path)
					assert.Equal(*WorkflowsBatches[iwf].OnExit[j].Content, *test.expectedWorkflowBatch[iwf].OnExit[j].Content)
				}

				for k, _ := range wf.Templates {
					if test.expectedWorkflowBatch[iwf].Templates[k].Path == nil || test.expectedWorkflowBatch[iwf].Templates[k].Content == nil {
						assert.Nil(WorkflowsBatches[iwf].Templates[k].Path)
						assert.Nil(WorkflowsBatches[iwf].Templates[k].Content)
					} else {
						assert.Equal(*WorkflowsBatches[iwf].Templates[k].Path, *test.expectedWorkflowBatch[iwf].Templates[k].Path)
						assert.Equal(*WorkflowsBatches[iwf].Templates[k].Content, *test.expectedWorkflowBatch[iwf].Templates[k].Content)
					}

				}

				if test.expectedWorkflowBatch[iwf].Parameters.Path == nil || test.expectedWorkflowBatch[iwf].Parameters.Content == nil {
					assert.Nil(WorkflowsBatches[iwf].Parameters.Path)
					assert.Nil(WorkflowsBatches[iwf].Parameters.Content)
				} else {
					assert.Equal(*WorkflowsBatches[iwf].Parameters.Path, *test.expectedWorkflowBatch[iwf].Parameters.Path)
				}
				assert.Equal(*WorkflowsBatches[iwf].Config, *test.expectedWorkflowBatch[iwf].Config)

			}
		})
	}

}
