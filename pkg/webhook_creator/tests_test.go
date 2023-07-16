package webhook_creator

import (
	"errors"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/git_provider"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"math/rand"
	"testing"
	"time"
)

func TestWebhookCreatorImpl_GetWebhook(t *testing.T) {
	assertion := assert.New(t)
	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{}, &clients.Clients{})

	// Add a webhook for testing
	hookID1 := int64(1)
	repoName1 := "test1"
	wc.setWebhook(hookID1, true, repoName1)

	hookID2 := int64(2)
	repoName2 := "test2"
	wc.setWebhook(hookID2, true, repoName2)

	// Retrieve the webhook
	hook := wc.getWebhook(hookID1)

	// Verify that the retrieved webhook matches the added webhook
	assertion.NotNil(hook)
	assertion.Equal(hookID1, hook.HookID)
	assertion.Equal(true, hook.HealthStatus)
	assertion.Equal(repoName1, *hook.RepoName)

	// Retrieve the webhook
	hook = wc.getWebhook(3)

	// Verify that the retrieved webhook matches the added webhook
	assertion.Nil(hook)

}

func TestWebhookCreatorImpl_SetWebhook(t *testing.T) {
	assertion := assert.New(t)
	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{}, &clients.Clients{})

	// Set a webhook
	hookID := int64(1)
	repoName := "test1"
	wc.setWebhook(hookID, true, repoName)

	// Get the webhook and verify the values
	hook := wc.getWebhook(hookID)

	assertion.NotNil(hook, "Webhook not set")
	assertion.Equal(hookID, hook.HookID, "Webhook HookID incorrect")
	assertion.True(hook.HealthStatus, "Webhook HealthStatus incorrect")
	assertion.Equal(repoName, *hook.RepoName, "Webhook RepoName incorrect")

}

func TestWebhookCreatorImpl_DeleteWebhook(t *testing.T) {
	assertion := assert.New(t)
	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{}, &clients.Clients{})

	// Add a webhook for testing
	hookID := int64(1)
	repoName := "test1"
	wc.setWebhook(hookID, true, repoName)

	// Delete the webhook
	wc.deleteWebhook(hookID)

	// Verify that the webhook was deleted
	hook := wc.getWebhook(hookID)
	assertion.Nil(hook)
}

func TestWebhookCreatorImpl_SetWebhookHealth(t *testing.T) {
	assertion := assert.New(t)
	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{}, &clients.Clients{})

	// Add a webhook for testing
	hookID := int64(1)
	repoName := "test1"
	wc.setWebhook(hookID, true, repoName)

	// Set the webhook health status
	err := wc.SetWebhookHealth(hookID, false)
	assertion.Nil(err, "SetWebhookHealth returned an error: %v", err)

	// Verify that the webhook health status was updated
	hook := wc.getWebhook(hookID)
	assertion.Equal(false, hook.HealthStatus, "Webhook health status not updated correctly")

	// Set the webhook health status of non existing webhook
	err = wc.SetWebhookHealth(2, false)
	assertion.NotNil(err)

}

func TestWebhookCreatorImpl_SetAllHooksHealth(t *testing.T) {
	assertion := assert.New(t)
	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{}, &clients.Clients{})

	// Add webhooks for testing
	wc.setWebhook(1, true, "repo1")
	wc.setWebhook(2, true, "repo2")
	wc.setWebhook(3, true, "repo3")

	// Set the health status of all webhooks to false
	wc.setAllHooksHealth(false)

	// Verify that all webhooks have the updated health status
	for _, hook := range wc.hooks {
		assertion.False(hook.HealthStatus, "Webhook health status not updated correctly")
	}
}

func TestWebhookCreatorImpl_InitWebhooks(t *testing.T) {
	assertion := assert.New(t)
	// Mock the necessary methods of the GitProvider client
	mockClient := &MockGitProviderClient{
		SetWebhookFunc: func(ctx context.Context, repoName *string) (*git_provider.HookWithStatus, error) {
			// Simulate setting a new webhook
			return &git_provider.HookWithStatus{
				HookID:       rand.Int63(),
				HealthStatus: true,
				RepoName:     repoName,
			}, nil
		},
	}

	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{
		GitProviderConfig: conf.GitProviderConfig{
			OrgLevelWebhook: false,
			RepoList:        "repo1,repo2",
		},
	}, &clients.Clients{
		GitProvider: mockClient,
	})

	// Initialize the webhooks
	err := wc.initWebhooks()

	// Run tests
	assertion.NoError(err)
	assertion.Len(wc.hooks, 2)
}

func TestWebhookCreatorImpl_Stop(t *testing.T) {
	assertion := assert.New(t)

	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{}, &clients.Clients{})

	// Add a webhook for testing
	hookID := int64(1)
	repoName := "repo1"
	wc.setWebhook(hookID, true, repoName)

	// Mock the necessary methods of the GitProvider client
	mockClient := &MockGitProviderClient{
		UnsetWebhookFunc: func(ctx context.Context, hook *git_provider.HookWithStatus) error {
			// Simulate successful webhook deletion
			return nil
		},
	}

	// Replace the GitProvider client in the WebhookCreatorImpl with the mock client
	wc.clients.GitProvider = mockClient

	// Stop the webhook creator
	ctx := context.Background()
	wc.Stop(&ctx)

	// Verify that the webhook was deleted
	hook := wc.getWebhook(hookID)
	assertion.Nil(hook, "Webhook not deleted")

}

func TestWebhookCreatorImpl_DeleteWebhooks(t *testing.T) {
	assertion := assert.New(t)

	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{}, &clients.Clients{})

	// Add webhooks for testing
	wc.setWebhook(1, true, "repo1")
	wc.setWebhook(2, true, "repo2")
	wc.setWebhook(3, true, "repo3")

	// Mock the necessary methods of the GitProvider client
	mockClient := &MockGitProviderClient{
		UnsetWebhookFunc: func(ctx context.Context, hook *git_provider.HookWithStatus) error {
			// Simulate successful webhook deletion
			return nil
		},
	}

	// Replace the GitProvider client in the WebhookCreatorImpl with the mock client
	wc.clients.GitProvider = mockClient

	// Delete the webhooks
	ctx := context.Background()
	err := wc.deleteWebhooks(&ctx)
	assertion.NoError(err)

	// Verify that all webhooks were deleted
	assertion.Len(wc.hooks, 0, "Webhooks not deleted")

}

func TestWebhookCreatorImpl_CheckHooksHealth(t *testing.T) {
	assertion := assert.New(t)

	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{}, &clients.Clients{})

	// Add webhooks for testing
	wc.setWebhook(1, true, "repo1")
	wc.setWebhook(2, false, "repo2")
	wc.setWebhook(3, true, "repo3")

	// Check the health sta tus of webhooks with a timeout of 1 second
	timeout := 1 * time.Second
	result := wc.checkHooksHealth(timeout)

	// Verify that the result indicates that all webhooks are healthy
	assertion.False(result)

	wc.setWebhook(2, true, "repo2")

	timeStart := time.Now()
	result = wc.checkHooksHealth(timeout)
	timeTook := time.Since(timeStart)
	assertion.Less(timeTook, timeout)
	assertion.True(result)
}

func TestWebhookCreatorImpl_RecoverHook(t *testing.T) {
	assertion := assert.New(t)

	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{}, &clients.Clients{})

	// Add a webhook for testing
	hookID := int64(1)
	repoName := "repo"
	wc.setWebhook(hookID, false, repoName)

	// Mock the necessary methods of the GitProvider client
	mockClient := &MockGitProviderClient{
		SetWebhookFunc: func(ctx context.Context, repoName *string) (*git_provider.HookWithStatus, error) {
			// Simulate setting a new webhook
			return &git_provider.HookWithStatus{
				HookID:       hookID,
				HealthStatus: true,
				RepoName:     repoName,
			}, nil
		},
	}

	// Replace the GitProvider client in the WebhookCreatorImpl with the mock client
	wc.clients.GitProvider = mockClient

	// Recover the webhook
	ctx := context.Background()
	err := wc.recoverHook(&ctx, hookID)
	assertion.Nil(err)

	// Verify that the webhook was recovered
	hook := wc.getWebhook(hookID)
	assertion.NotNil(hook)
	assertion.True(hook.HealthStatus)
}

func TestWebhookCreatorImpl_PingHooks(t *testing.T) {
	assertion := assert.New(t)

	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{
		GitProviderConfig: conf.GitProviderConfig{
			OrgLevelWebhook: true,
			RepoList:        "repo1,repo2,repo3",
		},
	}, &clients.Clients{})

	// Add webhooks for testing
	wc.setWebhook(1, true, "repo1")
	wc.setWebhook(2, false, "repo2")
	wc.setWebhook(3, true, "repo3")

	// Mock the necessary methods of the GitProvider client
	mockClient := &MockGitProviderClient{
		PingHookFunc: func(ctx context.Context, hook *git_provider.HookWithStatus) error {
			if hook.HookID == 2 {
				// Simulate a failure when pinging a specific webhook
				return errors.New("ping failed")
			}
			return nil
		},
		SetWebhookFunc: func(ctx context.Context, repoName *string) (*git_provider.HookWithStatus, error) {
			// Simulate setting a new webhook
			return &git_provider.HookWithStatus{
				HookID:       4,
				HealthStatus: true,
				RepoName:     repoName,
			}, nil
		},
	}

	// Replace the GitProvider client in the WebhookCreatorImpl with the mock client
	wc.clients.GitProvider = mockClient

	// Ping the webhooks
	ctx := context.Background()
	err := wc.pingHooks(&ctx)
	assertion.Nil(err)
	// Verify that the failed webhook was recovered
	hook := wc.getWebhook(4)
	assertion.NotNil(hook)
	assertion.True(hook.HealthStatus)
}

func TestWebhookCreatorImpl_RunDiagnosis(t *testing.T) {
	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{}, &clients.Clients{})

	// Add a webhook for testing
	hookID := int64(1)
	repoName := "example/repo"
	wc.setWebhook(hookID, false, repoName)

	// Mock the necessary methods of the GitProvider client
	mockClient := &MockGitProviderClient{
		PingHookFunc: func(ctx context.Context, hook *git_provider.HookWithStatus) error {
			// Simulate a failure when pinging the webhook
			return errors.New("ping failed")
		},
		SetWebhookFunc: func(ctx context.Context, repoName *string) (*git_provider.HookWithStatus, error) {
			// Simulate setting a new webhook
			return &git_provider.HookWithStatus{
				HookID:       hookID,
				HealthStatus: true,
				RepoName:     repoName,
			}, nil
		},
	}

	// Replace the GitProvider client in the WebhookCreatorImpl with the mock client
	wc.clients.GitProvider = mockClient

	// Run the webhook diagnosis
	ctx := context.Background()
	err := wc.RunDiagnosis(&ctx)
	if err != nil {
		t.Errorf("RunDiagnosis returned an error: %v", err)
	}

	// Verify that the webhook health status was updated after recovery
	hook := wc.getWebhook(hookID)
	if !hook.HealthStatus {
		t.Errorf("Webhook health status not updated after recovery")
	}
}
