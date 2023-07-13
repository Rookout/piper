package git_provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-github/v52/github"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/utils"
	assertion "github.com/stretchr/testify/assert"
)

func TestListFiles(t *testing.T) {
	// Prepare
	client, mux, _, teardown := setup()
	defer teardown()

	repoContent := &github.RepositoryContent{
		Type: utils.SPtr("file"),
		Name: utils.SPtr("exit.yaml"),
		Path: utils.SPtr(".workflows/exit.yaml"),
	}

	repoContent2 := &github.RepositoryContent{
		Type: utils.SPtr("file"),
		Name: utils.SPtr("main.yaml"),
		Path: utils.SPtr(".workflows/main.yaml"),
	}

	jsonBytes, _ := json.Marshal([]github.RepositoryContent{*repoContent, *repoContent2})

	mux.HandleFunc("/repos/test/test-repo1/contents/.workflows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		//testFormValues(t, r, values{})

		// Get the ref value from the URL query parameters
		ref := r.URL.Query().Get("ref")
		expectedRef := "branch1"

		// Check if the ref value matches the expected value
		if ref != expectedRef {
			http.Error(w, "Invalid ref value", http.StatusBadRequest)
			return
		}

		_, _ = fmt.Fprint(w, string(jsonBytes))
	})

	c := GithubClientImpl{
		client: client,
		cfg: &conf.GlobalConfig{
			GitProviderConfig: conf.GitProviderConfig{
				OrgLevelWebhook: false,
				OrgName:         "test",
				RepoList:        "test-repo1",
			},
		},
	}
	ctx := context.Background()

	// Execute
	actualContent, err := c.ListFiles(&ctx, "test-repo1", "branch1", ".workflows")
	expectedContent := []string{"exit.yaml", "main.yaml"}

	// Assert
	assert := assertion.New(t)
	assert.NotNil(t, err)
	assert.Equal(expectedContent, actualContent)

}

func TestSetStatus(t *testing.T) {
	// Prepare
	ctx := context.Background()
	assert := assertion.New(t)
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/test/test-repo1/statuses/test-commit", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{})

		w.WriteHeader(http.StatusCreated)
		jsonBytes := []byte(`{"status": "ok"}`)
		_, _ = fmt.Fprint(w, string(jsonBytes))
	})

	c := GithubClientImpl{
		client: client,
		cfg: &conf.GlobalConfig{
			GitProviderConfig: conf.GitProviderConfig{
				OrgLevelWebhook: false,
				OrgName:         "test",
				RepoList:        "test-repo1",
			},
		},
	}

	// Define test cases
	tests := []struct {
		name        string
		repo        *string
		commit      *string
		linkURL     *string
		status      *string
		message     *string
		wantedError error
	}{
		{
			name:        "Notify success",
			repo:        utils.SPtr("test-repo1"),
			commit:      utils.SPtr("test-commit"),
			linkURL:     utils.SPtr("https://argo"),
			status:      utils.SPtr("success"),
			message:     utils.SPtr(""),
			wantedError: nil,
		},
		{
			name:        "Notify pending",
			repo:        utils.SPtr("test-repo1"),
			commit:      utils.SPtr("test-commit"),
			linkURL:     utils.SPtr("https://argo"),
			status:      utils.SPtr("pending"),
			message:     utils.SPtr(""),
			wantedError: nil,
		},
		{
			name:        "Notify error",
			repo:        utils.SPtr("test-repo1"),
			commit:      utils.SPtr("test-commit"),
			linkURL:     utils.SPtr("https://argo"),
			status:      utils.SPtr("error"),
			message:     utils.SPtr("some message"),
			wantedError: nil,
		},
		{
			name:        "Notify failure",
			repo:        utils.SPtr("test-repo1"),
			commit:      utils.SPtr("test-commit"),
			linkURL:     utils.SPtr("https://argo"),
			status:      utils.SPtr("failure"),
			message:     utils.SPtr(""),
			wantedError: nil,
		},
		{
			name:        "Non managed repo",
			repo:        utils.SPtr("non-existing-repo"),
			commit:      utils.SPtr("test-commit"),
			linkURL:     utils.SPtr("https://argo"),
			status:      utils.SPtr("error"),
			message:     utils.SPtr(""),
			wantedError: errors.New("some error"),
		},
		{
			name:        "Non existing commit",
			repo:        utils.SPtr("test-repo1"),
			commit:      utils.SPtr("not-exists"),
			linkURL:     utils.SPtr("https://argo"),
			status:      utils.SPtr("error"),
			message:     utils.SPtr(""),
			wantedError: errors.New("some error"),
		},
		{
			name:        "Wrong URL",
			repo:        utils.SPtr("test-repo1"),
			commit:      utils.SPtr("test-commit"),
			linkURL:     utils.SPtr("argo"),
			status:      utils.SPtr("error"),
			message:     utils.SPtr(""),
			wantedError: errors.New("some error"),
		},
	}
	// Run test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// Call the function being tested
			err := c.SetStatus(&ctx, test.repo, test.commit, test.linkURL, test.status, test.message)

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

func TestSetWebhook(t *testing.T) {
	// Prepare
	ctx := context.Background()
	assert := assertion.New(t)
	client, mux, _, teardown := setup()
	defer teardown()

	hookUrl := "https://url"
	hooksList := []*github.Hook{
		&github.Hook{
			ID:     utils.IPtr(123),
			Name:   utils.SPtr("web"),
			Active: utils.BPtr(true),
			Events: []string{"pull_request", "create", "push"},
			Config: map[string]interface{}{
				"url": hookUrl,
			},
		},
	}

	// Existing webhook org
	mux.HandleFunc("/orgs/test/hooks", func(w http.ResponseWriter, r *http.Request) {
		var jsonBytes []byte
		if r.Method == "POST" {
			testFormValues(t, r, values{})
			w.WriteHeader(http.StatusCreated)
			jsonBytes, _ = json.Marshal(hooksList[0])
		}

		if r.Method == "GET" {
			testFormValues(t, r, values{})
			w.WriteHeader(http.StatusOK)
			jsonBytes, _ = json.Marshal(hooksList)
		}

		_, _ = fmt.Fprint(w, string(jsonBytes))
	})

	mux.HandleFunc("/orgs/test/hooks/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, values{})
		w.WriteHeader(http.StatusOK)
		jsonBytes, _ := json.Marshal(hooksList[0])
		_, _ = fmt.Fprint(w, string(jsonBytes))
	})

	// Not existing webhook org
	mux.HandleFunc("/orgs/test2/hooks", func(w http.ResponseWriter, r *http.Request) {
		var jsonBytes []byte
		if r.Method == "POST" {
			testFormValues(t, r, values{})
			w.WriteHeader(http.StatusCreated)
			jsonBytes, _ = json.Marshal(hooksList[0])
		}

		if r.Method == "GET" {
			testFormValues(t, r, values{})
			w.WriteHeader(http.StatusOK)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		_, _ = fmt.Fprint(w, string(jsonBytes))
	})

	mux.HandleFunc("/orgs/test2/hooks/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, values{})
		w.WriteHeader(http.StatusOK)
		jsonBytes, _ := json.Marshal(hooksList[0])
		_, _ = fmt.Fprint(w, string(jsonBytes))
	})

	// Existing webhook repo
	mux.HandleFunc("/repos/test/test-repo1/hooks", func(w http.ResponseWriter, r *http.Request) {
		var jsonBytes []byte
		if r.Method == "POST" {
			testFormValues(t, r, values{})
			w.WriteHeader(http.StatusCreated)
			jsonBytes, _ = json.Marshal(hooksList[0])
		}

		if r.Method == "GET" {
			testFormValues(t, r, values{})
			w.WriteHeader(http.StatusOK)
			jsonBytes, _ = json.Marshal(hooksList)
		}

		_, _ = fmt.Fprint(w, string(jsonBytes))
	})

	mux.HandleFunc("/repos/test/test-repo1/hooks/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, values{})
		w.WriteHeader(http.StatusOK)
		jsonBytes, _ := json.Marshal(hooksList[0])
		_, _ = fmt.Fprint(w, string(jsonBytes))
	})

	// Not existing webhook repo
	mux.HandleFunc("/repos/test/test-repo2/hooks/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, values{})
		w.WriteHeader(http.StatusOK)
		jsonBytes, _ := json.Marshal(hooksList[0])
		_, _ = fmt.Fprint(w, string(jsonBytes))
	})

	mux.HandleFunc("/repos/test/test-repo2/hooks", func(w http.ResponseWriter, r *http.Request) {
		var jsonBytes []byte
		if r.Method == "POST" {
			testFormValues(t, r, values{})
			w.WriteHeader(http.StatusCreated)
			jsonBytes, _ = json.Marshal(hooksList[0])
		}

		if r.Method == "GET" {
			testFormValues(t, r, values{})
			w.WriteHeader(http.StatusNotFound)
			return
		}

		_, _ = fmt.Fprint(w, string(jsonBytes))
	})

	c := GithubClientImpl{
		client: client,
		cfg: &conf.GlobalConfig{
			GitProviderConfig: conf.GitProviderConfig{},
		},
	}

	// Define test cases
	tests := []struct {
		name        string
		repo        *string
		config      *conf.GitProviderConfig
		wantedError error
	}{
		{
			name: "Set repo webhook",
			repo: utils.SPtr("test-repo1"),
			config: &conf.GitProviderConfig{
				OrgLevelWebhook: false,
				OrgName:         "test",
				RepoList:        "test-repo1",
				WebhookURL:      hookUrl,
			},
			wantedError: nil,
		},
		{
			name: "Create repo webhook",
			repo: utils.SPtr("test-repo1"),
			config: &conf.GitProviderConfig{
				OrgLevelWebhook: false,
				OrgName:         "test",
				RepoList:        "test-repo2",
				WebhookURL:      hookUrl,
			},
			wantedError: nil,
		},
		{
			name: "Set org webhook",
			repo: nil,
			config: &conf.GitProviderConfig{
				OrgLevelWebhook: true,
				OrgName:         "test",
				RepoList:        "",
				WebhookURL:      hookUrl,
				WebhookSecret:   "test-secret",
			},
			wantedError: nil,
		},
		{
			name: "Create org webhook",
			repo: nil,
			config: &conf.GitProviderConfig{
				OrgLevelWebhook: true,
				OrgName:         "test2",
				RepoList:        "",
				WebhookURL:      hookUrl,
				WebhookSecret:   "test-secret",
			},
			wantedError: nil,
		},
		{
			name: "Set org with given repo",
			repo: utils.SPtr("test-repo1"),
			config: &conf.GitProviderConfig{
				OrgLevelWebhook: true,
				OrgName:         "test",
				RepoList:        "test-repo1",
				WebhookURL:      hookUrl,
				WebhookSecret:   "test-secret",
			},
			wantedError: errors.New("some error"),
		},
	}
	// Run test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c.cfg.GitProviderConfig = *test.config
			// Call the function being tested
			_, err := c.SetWebhook(&ctx, test.repo)

			// Use assert to check the equality of the error
			if test.wantedError != nil {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				//assert.Equal(hookUrl, hook.Config["url"])
			}
		})
	}

}

//
//func TestPingHook(t *testing.T) {
//	// Prepare
//	ctx := context.Background()
//	assert := assertion.New(t)
//	client, mux, _, teardown := setup()
//	defer teardown()
//
//	//hookUrl := "https://url"
//	orgHooksList := []*HookWithStatus{
//		{
//			HookID: utils.IPtr(123),
//			//Hook: &github.Hook{
//			//	ID:     utils.IPtr(123),
//			//	Name:   utils.SPtr("web"),
//			//	Active: utils.BPtr(true),
//			//	Events: []string{"pull_request", "create", "push"},
//			//	Config: map[string]interface{}{
//			//		"url": hookUrl,
//			//	},
//			//},
//			HealthStatus: true,
//			RepoName:     nil,
//		},
//	}
//
//	repoHooksList := []*HookWithStatus{
//		{
//			HookID: utils.IPtr(234),
//			//Hook: &github.Hook{
//			//	ID:     utils.IPtr(234),
//			//	Name:   utils.SPtr("web"),
//			//	Active: utils.BPtr(true),
//			//	Events: []string{"pull_request", "create", "push"},
//			//	Config: map[string]interface{}{
//			//		"url": hookUrl,
//			//	},
//			//},
//			HealthStatus: true,
//			RepoName:     utils.SPtr("test-repo1"),
//		},
//	}
//	// Test-repo2 existing webhook
//	mux.HandleFunc("/repos/test/test-repo1/hooks/234/pings", func(w http.ResponseWriter, r *http.Request) {
//		testMethod(t, r, "POST")
//		testFormValues(t, r, values{})
//		w.WriteHeader(http.StatusNoContent)
//	})
//
//	mux.HandleFunc("/orgs/test/hooks/123/pings", func(w http.ResponseWriter, r *http.Request) {
//		testMethod(t, r, "POST")
//		testFormValues(t, r, values{})
//		w.WriteHeader(http.StatusNoContent)
//	})
//
//	c := GithubClientImpl{
//		client: client,
//		cfg: &conf.GlobalConfig{
//			GitProviderConfig: conf.GitProviderConfig{},
//		},
//	}
//
//	// Define test cases
//	tests := []struct {
//		name        string
//		repo        *string
//		hooks       []*HookWithStatus
//		config      *conf.GitProviderConfig
//		wantedError error
//	}{
//		{
//			name:  "Ping repo webhook",
//			hooks: repoHooksList,
//			config: &conf.GitProviderConfig{
//				OrgLevelWebhook: false,
//				OrgName:         "test",
//			},
//			wantedError: nil,
//		},
//		{
//			name:  "Ping org webhook",
//			hooks: orgHooksList,
//			config: &conf.GitProviderConfig{
//				OrgLevelWebhook: true,
//				OrgName:         "test",
//			},
//			wantedError: nil,
//		},
//	}
//	// Run test cases
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			c.hooks = test.hooks
//			c.cfg.GitProviderConfig = *test.config
//			// Call the function being tested
//			err := c.PingHooks(&ctx)
//
//			// Use assert to check the equality of the error
//			if test.wantedError != nil {
//				assert.NotNil(err)
//			} else {
//				assert.Nil(err)
//			}
//		})
//	}
//
//}
