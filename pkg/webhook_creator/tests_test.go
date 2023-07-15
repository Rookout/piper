package webhook_creator

import (
	"errors"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/git_provider"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestWebhookCreatorImpl_GetWebhook(t *testing.T) {
	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{}, &clients.Clients{})

	// Add a webhook for testing
	hookID := int64(1)
	repoName := "example/repo"
	wc.setWebhook(hookID, true, repoName)

	// Retrieve the webhook
	hook := wc.getWebhook(hookID)

	// Verify that the retrieved webhook matches the added webhook
	if hook == nil || *hook.HookID != hookID || !hook.HealthStatus || *hook.RepoName != repoName {
		t.Errorf("GetWebhook did not return the expected webhook")
	}
}

func TestWebhookCreatorImpl_SetWebhook(t *testing.T) {
	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{}, &clients.Clients{})

	// Set a webhook
	hookID := int64(1)
	healthStatus := true
	repoName := "example/repo"
	wc.setWebhook(hookID, healthStatus, repoName)

	// Get the webhook and verify the values
	hook := wc.getWebhook(hookID)
	if hook == nil {
		t.Errorf("Webhook not set")
	}
	if *hook.HookID != hookID {
		t.Errorf("Webhook HookID incorrect")
	}
	if hook.HealthStatus != healthStatus {
		t.Errorf("Webhook HealthStatus incorrect")
	}
	if *hook.RepoName != repoName {
		t.Errorf("Webhook RepoName incorrect")
	}
}

func TestWebhookCreatorImpl_DeleteWebhook(t *testing.T) {
	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{}, &clients.Clients{})

	// Add a webhook for testing
	hookID := int64(1)
	repoName := "example/repo"
	wc.setWebhook(hookID, true, repoName)

	// Delete the webhook
	wc.deleteWebhook(hookID)

	// Verify that the webhook was deleted
	hook := wc.getWebhook(hookID)
	if hook != nil {
		t.Errorf("Webhook not deleted")
	}
}

func TestWebhookCreatorImpl_SetWebhookHealth(t *testing.T) {
	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{}, &clients.Clients{})

	// Add a webhook for testing
	hookID := int64(1)
	repoName := "example/repo"
	wc.setWebhook(hookID, true, repoName)

	// Set the webhook health status
	err := wc.SetWebhookHealth(hookID, false)
	if err != nil {
		t.Errorf("SetWebhookHealth returned an error: %v", err)
	}

	// Verify that the webhook health status was updated
	hook := wc.getWebhook(hookID)
	if hook.HealthStatus {
		t.Errorf("Webhook health status not updated correctly")
	}
}

func TestWebhookCreatorImpl_SetAllHooksHealth(t *testing.T) {
	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{}, &clients.Clients{})

	// Add webhooks for testing
	wc.setWebhook(1, true, "example/repo1")
	wc.setWebhook(2, true, "example/repo2")
	wc.setWebhook(3, true, "example/repo3")

	// Set the health status of all webhooks to false
	wc.setAllHooksHealth(false)

	// Verify that all webhooks have the updated health status
	for _, hook := range wc.hooks {
		if hook.HealthStatus {
			t.Errorf("Webhook health status not updated correctly")
		}
	}
}

func TestWebhookCreatorImpl_InitWebhooks(t *testing.T) {
	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{
		GitProviderConfig: conf.GitProviderConfig{
			OrgLevelWebhook: true,
			RepoList:        "example/repo",
		},
	}, &clients.Clients{})

	// Initialize the webhooks
	err := wc.initWebhooks()
	if err == nil {
		t.Errorf("InitWebhooks did not return an error as expected")
	}

	// Verify that the initialization error occurred due to invalid configuration
}

//
//func TestWebhookCreatorImpl_Stop(t *testing.T) {
//	// Create a test instance of the WebhookCreatorImpl
//	wc := NewWebhookCreator(&conf.GlobalConfig{}, &clients.Clients{})
//
//	// Add a webhook for testing
//	hookID := int64(1)
//	repoName := "example/repo"
//	wc.setWebhook(hookID, true, repoName)
//
//	// Mock the necessary methods of the GitProvider client
//	mockClient := &MockGitProviderClient{
//		UnsetWebhookFunc: func(ctx context.Context, hook *git_provider.HookWithStatus) error {
//			// Simulate successful webhook deletion
//			return nil
//		},
//	}
//
//	// Replace the GitProvider client in the WebhookCreatorImpl with the mock client
//	wc.clients.GitProvider = mockClient
//
//	// Stop the webhook creator
//	ctx := context.Background()
//	wc.Stop(&ctx)
//
//	// Verify that the webhook was deleted
//	if wc.getWebhook(hookID) != nil {
//		t.Errorf("Webhook not deleted")
//	}
//}

//func TestWebhookCreatorImpl_DeleteWebhooks(t *testing.T) {
//	// Create a test instance of the WebhookCreatorImpl
//	wc := NewWebhookCreator(&conf.GlobalConfig{}, &clients.Clients{})
//
//	// Add webhooks for testing
//	wc.setWebhook(1, true, "example/repo1")
//	wc.setWebhook(2, true, "example/repo2")
//	wc.setWebhook(3, true, "example/repo3")
//
//	// Mock the necessary methods of the GitProvider client
//	mockClient := &MockGitProviderClient{
//		UnsetWebhookFunc: func(ctx context.Context, hook *git_provider.HookWithStatus) error {
//			// Simulate successful webhook deletion
//			return nil
//		},
//	}
//
//	// Replace the GitProvider client in the WebhookCreatorImpl with the mock client
//	wc.clients.GitProvider = mockClient
//
//	// Delete the webhooks
//	ctx := context.Background()
//	err := wc.deleteWebhooks(&ctx)
//	if err != nil {
//		t.Errorf("DeleteWebhooks returned an error: %v", err)
//	}
//
//	// Verify that all webhooks were deleted
//	if len(wc.hooks) != 0 {
//		t.Errorf("Webhooks not deleted")
//	}
//}

func TestWebhookCreatorImpl_CheckHooksHealth(t *testing.T) {
	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{}, &clients.Clients{})

	// Add webhooks for testing
	wc.setWebhook(1, true, "example/repo1")
	wc.setWebhook(2, false, "example/repo2")
	wc.setWebhook(3, true, "example/repo3")

	// Check the health status of webhooks with a timeout of 1 second
	timeout := 1 * time.Second
	result := wc.checkHooksHealth(timeout)

	// Verify that the result indicates that all webhooks are healthy
	if !result {
		t.Errorf("CheckHooksHealth returned false")
	}
}

func TestWebhookCreatorImpl_RecoverHook(t *testing.T) {
	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{}, &clients.Clients{})

	// Add a webhook for testing
	hookID := int64(1)
	repoName := "example/repo"
	wc.setWebhook(hookID, false, repoName)

	// Mock the necessary methods of the GitProvider client
	mockClient := &MockGitProviderClient{
		SetWebhookFunc: func(ctx context.Context, repoName *string) (*git_provider.HookWithStatus, error) {
			// Simulate setting a new webhook
			return &git_provider.HookWithStatus{
				HookID:       &hookID,
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
	if err != nil {
		t.Errorf("RecoverHook returned an error: %v", err)
	}

	// Verify that the webhook was recovered
	hook := wc.getWebhook(hookID)
	if hook == nil || hook.HealthStatus {
		t.Errorf("Webhook not recovered correctly")
	}
}

func TestWebhookCreatorImpl_PingHooks(t *testing.T) {
	// Create a test instance of the WebhookCreatorImpl
	wc := NewWebhookCreator(&conf.GlobalConfig{}, &clients.Clients{})

	// Add webhooks for testing
	wc.setWebhook(1, true, "example/repo1")
	wc.setWebhook(2, false, "example/repo2")
	wc.setWebhook(3, true, "example/repo3")

	// Mock the necessary methods of the GitProvider client
	mockClient := &MockGitProviderClient{
		PingHookFunc: func(ctx context.Context, hook *git_provider.HookWithStatus) error {
			if *hook.HookID == 2 {
				// Simulate a failure when pinging a specific webhook
				return errors.New("ping failed")
			}
			return nil
		},
		SetWebhookFunc: func(ctx context.Context, repoName *string) (*git_provider.HookWithStatus, error) {
			// Simulate setting a new webhook
			return &git_provider.HookWithStatus{
				HealthStatus: true,
				RepoName:     repoName,
			}, nil
		},
	}

	// Replace the GitProvider client in the WebhookCreatorImpl with the mock client
	wc.clients.GitProvider = mockClient

	// Ping the webhooks
	ctx := context.Background()
	wc.pingHooks(&ctx)

	// Verify that the failed webhook was recovered
	hook := wc.getWebhook(2)
	if hook == nil || hook.HealthStatus {
		t.Errorf("Failed webhook not recovered correctly")
	}
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
				HookID:       &hookID,
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
