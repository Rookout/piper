package event_handler

import (
	"context"
	"errors"
	"github.com/rookout/piper/pkg/git_provider"
	assertion "github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"testing"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
)

type mockGitProvider struct{}

func (m *mockGitProvider) GetFile(ctx *context.Context, repo string, branch string, path string) (*git_provider.CommitFile, error) {
	return nil, nil
}

func (m *mockGitProvider) GetFiles(ctx *context.Context, repo string, branch string, paths []string) ([]*git_provider.CommitFile, error) {
	return nil, nil
}

func (m *mockGitProvider) ListFiles(ctx *context.Context, repo string, branch string, path string) ([]string, error) {
	return nil, nil
}

func (m *mockGitProvider) SetWebhook() error {
	return nil
}

func (m *mockGitProvider) UnsetWebhook(ctx *context.Context) error {
	return nil
}

func (m *mockGitProvider) HandlePayload(request *http.Request, secret []byte) (*git_provider.WebhookPayload, error) {
	return nil, nil
}

func (m *mockGitProvider) SetStatus(ctx *context.Context, repo *string, commit *string, linkURL *string, status *string, message *string) error {
	return nil
}

func TestNotify(t *testing.T) {
	assert := assertion.New(t)
	ctx := context.Background()

	// Define test cases
	tests := []struct {
		name        string
		workflow    *v1alpha1.Workflow
		wantedError error
	}{
		{
			name: "Succeeded workflow",
			workflow: &v1alpha1.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-workflow",
					Labels: map[string]string{
						"repo":   "test-repo",
						"commit": "test-commit",
					},
				},
				Status: v1alpha1.WorkflowStatus{
					Phase:   v1alpha1.WorkflowSucceeded,
					Message: "",
				},
			},
			wantedError: nil,
		},
		{
			name: "Failed workflow",
			workflow: &v1alpha1.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-workflow",
					Labels: map[string]string{
						"repo":   "test-repo",
						"commit": "test-commit",
					},
				},
				Status: v1alpha1.WorkflowStatus{
					Phase:   v1alpha1.WorkflowFailed,
					Message: "something",
				},
			},
			wantedError: nil,
		},
		{
			name: "Error workflow",
			workflow: &v1alpha1.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-workflow",
					Labels: map[string]string{
						"repo":   "test-repo",
						"commit": "test-commit",
					},
				},
				Status: v1alpha1.WorkflowStatus{
					Phase:   v1alpha1.WorkflowError,
					Message: "something",
				},
			},
			wantedError: nil,
		},
		{
			name: "Pending workflow",
			workflow: &v1alpha1.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-workflow",
					Labels: map[string]string{
						"repo":   "test-repo",
						"commit": "test-commit",
					},
				},
				Status: v1alpha1.WorkflowStatus{
					Phase:   v1alpha1.WorkflowPending,
					Message: "something",
				},
			},
			wantedError: nil,
		},
		{
			name: "Running workflow",
			workflow: &v1alpha1.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-workflow",
					Labels: map[string]string{
						"repo":   "test-repo",
						"commit": "test-commit",
					},
				},
				Status: v1alpha1.WorkflowStatus{
					Phase:   v1alpha1.WorkflowRunning,
					Message: "something",
				},
			},
			wantedError: nil,
		},
		{
			name: "Missing label repo",
			workflow: &v1alpha1.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-workflow",
					Labels: map[string]string{
						"commit": "test-commit",
					},
				},
				Status: v1alpha1.WorkflowStatus{
					Phase:   v1alpha1.WorkflowSucceeded,
					Message: "something",
				},
			},
			wantedError: errors.New("some error"),
		},
		{
			name: "Missing label commit",
			workflow: &v1alpha1.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-workflow",
					Labels: map[string]string{
						"repo": "test-repo",
					},
				},
				Status: v1alpha1.WorkflowStatus{
					Phase:   v1alpha1.WorkflowSucceeded,
					Message: "something",
				},
			},
			wantedError: errors.New("some error"),
		},
	}

	// Create a mock configuration and clients
	cfg := &conf.GlobalConfig{
		WorkflowServerConfig: conf.WorkflowServerConfig{
			ArgoAddress: "http://workflow-server",
			Namespace:   "test-namespace",
		},
	}
	clients := &clients.Clients{
		GitProvider: &mockGitProvider{},
	}

	// Create a new githubNotifier instance
	gn := NewGithubEventNotifier(cfg, clients)

	// Call the Notify method

	// Run test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Call the function being tested
			err := gn.Notify(&ctx, test.workflow)

			// Use assert to check the equality of the error
			if test.wantedError != nil {
				assert.Error(err)
				assert.NotNil(err)
			} else {
				assert.NoError(err)
				assert.Nil(err)
			}
		})
	}
}
