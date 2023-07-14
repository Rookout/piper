package webhook_creator

//
//import (
//	"fmt"
//	"github.com/rookout/piper/pkg/clients"
//	"github.com/rookout/piper/pkg/git_provider"
//	"github.com/rookout/piper/pkg/utils"
//	"github.com/stretchr/testify/assert"
//	"golang.org/x/net/context"
//	"net/http"
//	"testing"
//	"time"
//)
//
//type MockGitProvider struct {
//	PingHookFunc     func(ctx *context.Context, hook *git_provider.HookWithStatus) error
//	SetWebhookFunc   func(ctx *context.Context, repo *string) (*git_provider.HookWithStatus, error)
//	UnsetWebhookFunc func(ctx *context.Context, hook *git_provider.HookWithStatus) error
//}
//
//func (m *MockGitProvider) PingHook(ctx *context.Context, hook *git_provider.HookWithStatus) error {
//	if m.PingHookFunc != nil {
//		return m.PingHookFunc(ctx, hook)
//	}
//	return nil
//}
//
//func (m *MockGitProvider) SetWebhook(ctx *context.Context, repo *string) (*git_provider.HookWithStatus, error) {
//	if m.SetWebhookFunc != nil {
//		return m.SetWebhookFunc(ctx, repo)
//	}
//	return nil, nil
//}
//
//func (m *MockGitProvider) UnsetWebhook(ctx *context.Context, hook *git_provider.HookWithStatus) error {
//	if m.UnsetWebhookFunc != nil {
//		return m.UnsetWebhookFunc(ctx, hook)
//	}
//	return nil
//}
//
//func (m *MockGitProvider) GetFile(ctx *context.Context, repo string, branch string, path string) (*git_provider.CommitFile, error) {
//	return nil, nil
//}
//
//func (m *MockGitProvider) ListFiles(ctx *context.Context, repo string, branch string, path string) ([]string, error) {
//	return nil, nil
//}
//
//func (m *MockGitProvider) GetFiles(ctx *context.Context, repo string, branch string, paths []string) ([]*git_provider.CommitFile, error) {
//	return nil, nil
//}
//
//func (m *MockGitProvider) HandlePayload(request *http.Request, secret []byte) (*git_provider.WebhookPayload, error) {
//	return nil, nil
//}
//
//func (m *MockGitProvider) SetStatus(ctx *context.Context, repo *string, commit *string, linkURL *string, status *string, message *string) error {
//	return nil
//}
//
//func TestWebhookCreatorImpl_SetHealth(t *testing.T) {
//	assertion := assert.New(t)
//
//	wc := &WebhookCreatorImpl{
//		clients: &clients.Clients{
//			GitProvider: &MockGitProvider{},
//		},
//		hooks: make(map[int64]*git_provider.HookWithStatus),
//	}
//
//	hookID := int64(123)
//	wc.hooks[hookID] = &git_provider.HookWithStatus{
//		HookID:       &hookID,
//		HealthStatus: false,
//	}
//
//	// Test setting health status to true
//	err := wc.SetWebhookHealth(true, &hookID)
//	assertion.NoError(err, "SetWebhookHealth should not return an error")
//	assertion.True(wc.hooks[hookID].HealthStatus, "SetWebhookHealth did not update the health status to true")
//
//	// Test setting health status to false
//	err = wc.SetWebhookHealth(false, &hookID)
//	assertion.NoError(err, "SetWebhookHealth should not return an error")
//	assertion.False(wc.hooks[hookID].HealthStatus, "SetWebhookHealth did not update the health status to false")
//
//	// Test error when hookID is not found
//	err = wc.SetWebhookHealth(true, &hookID) // Using the same hookID
//	assertion.Error(err, "SetWebhookHealth should return an error when hookID is not found")
//}
//
//func TestWebhookCreatorImpl_RunDiagnosis(t *testing.T) {
//	assertion := assert.New(t)
//
//	wc := &WebhookCreatorImpl{
//		clients: &clients.Clients{
//			GitProvider: &MockGitProvider{},
//		},
//		hooks: make(map[int64]*git_provider.HookWithStatus),
//	}
//
//	// Add a healthy hook
//	hookID1 := int64(1)
//	wc.hooks[hookID1] = &git_provider.HookWithStatus{
//		HookID:       &hookID1,
//		HealthStatus: true,
//	}
//
//	// Add an unhealthy hook
//	hookID2 := int64(2)
//	wc.hooks[hookID2] = &git_provider.HookWithStatus{
//		HookID:       &hookID2,
//		HealthStatus: false,
//	}
//
//	ctx := context.Background()
//
//	err := wc.RunDiagnosis(&ctx)
//	assertion.NoError(err, "RunDiagnosis should not return an error")
//
//	// Check that the unhealthy hook has been recovered
//	assertion.True(wc.hooks[hookID2].HealthStatus, "RunDiagnosis did not recover the unhealthy hook")
//}
//
//func TestWebhookCreatorImpl_SetWebhooks(t *testing.T) {
//	assertion := assert.New(t)
//
//	wc := &WebhookCreatorImpl{
//		clients: &clients.Clients{
//			GitProvider: &MockGitProvider{},
//		},
//		hooks: make(map[int64]*git_provider.HookWithStatus),
//	}
//
//	repoList := "repo1,repo2,repo3"
//
//	// Test setting webhooks successfully
//	wc.cfg.GitProviderConfig.RepoList = repoList
//	err := wc.setWebhooks()
//	assertion.NoError(err, "setWebhooks should not return an error")
//
//	// Check that the correct number of hooks has been created
//	assertion.Len(wc.hooks, 3, "setWebhooks did not create the correct number of hooks")
//
//	// Test error when org level webhook wanted but repositories list is provided
//	wc.cfg.GitProviderConfig.OrgLevelWebhook = true
//	err = wc.setWebhooks()
//	assertion.Error(err, "setWebhooks should return an error when org level webhook wanted but repositories list is provided")
//}
//
//func TestWebhookCreatorImpl_UnsetWebhooks(t *testing.T) {
//	assertion := assert.New(t)
//
//	wc := &WebhookCreatorImpl{
//		clients: &clients.Clients{
//			GitProvider: &MockGitProvider{},
//		},
//		hooks: make(map[int64]*git_provider.HookWithStatus),
//	}
//
//	// Add some hooks
//	wc.hooks[1] = &git_provider.HookWithStatus{
//		HookID: func() *int64 { id := int64(1); return &id }(),
//	}
//	wc.hooks[2] = &git_provider.HookWithStatus{
//		HookID: func() *int64 { id := int64(2); return &id }(),
//	}
//
//	ctx := context.Background()
//
//	// Test unsetting webhooks successfully
//	err := wc.unsetWebhooks(&ctx)
//	assertion.NoError(err, "unsetWebhooks should not return an error")
//
//	// Check that all hooks have been removed
//	assertion.Empty(wc.hooks, "unsetWebhooks did not remove all hooks")
//}
//
//func TestWebhookCreatorImpl_CheckHooksHealth(t *testing.T) {
//	assertion := assert.New(t)
//
//	wc := &WebhookCreatorImpl{
//		clients: &clients.Clients{
//			GitProvider: &MockGitProvider{},
//		},
//		hooks: make(map[int64]*git_provider.HookWithStatus),
//	}
//
//	// Add healthy hooks
//	wc.hooks[1] = &git_provider.HookWithStatus{
//		HookID:       func() *int64 { id := int64(1); return &id }(),
//		HealthStatus: true,
//	}
//	wc.hooks[2] = &git_provider.HookWithStatus{
//		HookID:       func() *int64 { id := int64(2); return &id }(),
//		HealthStatus: true,
//	}
//
//	// Add an unhealthy hook
//	wc.hooks[3] = &git_provider.HookWithStatus{
//		HookID:       func() *int64 { id := int64(3); return &id }(),
//		HealthStatus: false,
//	}
//
//	timeout := 5 * time.Second
//
//	// Test when all hooks are healthy
//	allHealthy := wc.checkHooksHealth(timeout)
//	assertion.True(allHealthy, "checkHooksHealth should return true when all hooks are healthy")
//
//	// Test when a hook becomes unhealthy during the timeout period
//	go func() {
//		time.Sleep(2 * time.Second)
//		wc.hooks[2].HealthStatus = false
//	}()
//
//	allHealthy = wc.checkHooksHealth(timeout)
//	assertion.False(allHealthy, "checkHooksHealth should return false when a hook becomes unhealthy")
//
//	// Test when the timeout period expires
//	allHealthy = wc.checkHooksHealth(1 * time.Second)
//	assertion.False(allHealthy, "checkHooksHealth should return false when the timeout period expires")
//}
//
//func TestWebhookCreatorImpl_RecoverHook(t *testing.T) {
//	assertion := assert.New(t)
//
//	wc := &WebhookCreatorImpl{
//		clients: &clients.Clients{
//			GitProvider: &MockGitProvider{},
//		},
//		hooks: make(map[int64]*git_provider.HookWithStatus),
//	}
//
//	hookID := int64(1)
//	wc.hooks[hookID] = &git_provider.HookWithStatus{
//		HookID:   func() *int64 { return &hookID }(),
//		RepoName: utils.SPtr("example-repo"),
//	}
//
//	ctx := context.Background()
//
//	// Test recovering an existing hook
//	err := wc.recoverHook(&ctx, &hookID)
//	assertion.NoError(err, "recoverHook should not return an error")
//
//	// Test error when hookID is not found
//	err = wc.recoverHook(&ctx, func() *int64 { id := int64(999); return &id }())
//	assertion.Error(err, "recoverHook should return an error when hookID is not found")
//}
//
//func TestWebhookCreatorImpl_Stop(t *testing.T) {
//	assertion := assert.New(t)
//
//	wc := &WebhookCreatorImpl{
//		clients: &clients.Clients{
//			GitProvider: &MockGitProvider{},
//		},
//		hooks: make(map[int64]*git_provider.HookWithStatus),
//	}
//
//	// Add some hooks
//	wc.hooks[1] = &git_provider.HookWithStatus{
//		HookID: func() *int64 { id := int64(1); return &id }(),
//	}
//	wc.hooks[2] = &git_provider.HookWithStatus{
//		HookID: func() *int64 { id := int64(2); return &id }(),
//	}
//
//	ctx := context.Background()
//
//	// Test stopping and unsetting webhooks successfully
//	wc.Stop(&ctx)
//
//	// Check that all hooks have been removed
//	assertion.Empty(wc.hooks, "Stop should remove all hooks")
//}
//
//func TestWebhookCreatorImpl_SetAllHooksHealth(t *testing.T) {
//	assertion := assert.New(t)
//
//	wc := &WebhookCreatorImpl{
//		clients: &clients.Clients{
//			GitProvider: &MockGitProvider{},
//		},
//		hooks: make(map[int64]*git_provider.HookWithStatus),
//	}
//
//	// Add some hooks
//	wc.hooks[1] = &git_provider.HookWithStatus{
//		HookID:       func() *int64 { id := int64(1); return &id }(),
//		HealthStatus: false,
//	}
//	wc.hooks[2] = &git_provider.HookWithStatus{
//		HookID:       func() *int64 { id := int64(2); return &id }(),
//		HealthStatus: false,
//	}
//	wc.hooks[3] = &git_provider.HookWithStatus{
//		HookID:       func() *int64 { id := int64(3); return &id }(),
//		HealthStatus: false,
//	}
//
//	// Set all hooks to healthy
//	wc.setAllHooksHealth(true)
//
//	for _, hook := range wc.hooks {
//		assertion.True(hook.HealthStatus, "setAllHooksHealth did not set all hooks to healthy")
//	}
//
//	// Set all hooks to unhealthy
//	wc.setAllHooksHealth(false)
//
//	for _, hook := range wc.hooks {
//		assertion.False(hook.HealthStatus, "setAllHooksHealth did not set all hooks to unhealthy")
//	}
//}
//
//func TestWebhookCreatorImpl_PingHooks(t *testing.T) {
//	assertion := assert.New(t)
//
//	wc := &WebhookCreatorImpl{
//		clients: &clients.Clients{
//			GitProvider: &MockGitProvider{},
//		},
//		hooks: make(map[int64]*git_provider.HookWithStatus),
//	}
//
//	// Add some hooks
//	wc.hooks[1] = &git_provider.HookWithStatus{
//		HookID:       func() *int64 { id := int64(1); return &id }(),
//		HealthStatus: true,
//	}
//	wc.hooks[2] = &git_provider.HookWithStatus{
//		HookID:       func() *int64 { id := int64(2); return &id }(),
//		HealthStatus: false,
//	}
//	wc.hooks[3] = &git_provider.HookWithStatus{
//		HookID:       func() *int64 { id := int64(3); return &id }(),
//		HealthStatus: true,
//	}
//
//	ctx := context.Background()
//
//	// Test pinging all hooks successfully
//	wc.pingHooks(&ctx)
//
//	// Check that hooks remain unchanged
//	for _, hook := range wc.hooks {
//		assertion.Equal(true, hook.HealthStatus, "pingHooks should not change the health status of hooks")
//	}
//
//	// Test recovering an unhealthy hook
//	// Set a mock implementation for PingHook that returns an error
//	wc.clients.GitProvider = &MockGitProvider{
//		PingHookFunc: func(ctx *context.Context, hook *git_provider.HookWithStatus) error {
//			if *hook.HookID == 2 {
//				return fmt.Errorf("ping error")
//			}
//			return nil
//		},
//	}
//
//	wc.pingHooks(&ctx)
//
//	// Check that the unhealthy hook has been recovered
//	assertion.True(wc.hooks[2].HealthStatus, "pingHooks did not recover the unhealthy hook")
//}
