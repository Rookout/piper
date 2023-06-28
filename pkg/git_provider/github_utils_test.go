package git_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-github/v52/github"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/utils"
	assertion "github.com/stretchr/testify/assert"
)

func TestIsOrgWebhookEnabled(t *testing.T) {
	//
	// Prepare
	//
	client, mux, _, teardown := setup()
	defer teardown()

	config := make(map[string]interface{})
	config["url"] = "https://bla.com"
	Hooks := github.Hook{
		Active: utils.BPtr(true),
		Name:   utils.SPtr("web"),
		Config: config,
	}
	jsonBytes, _ := json.Marshal(&[]github.Hook{Hooks})

	mux.HandleFunc("/orgs/test/hooks", func(w http.ResponseWriter, r *http.Request) {
		TestMethod(t, r, "GET")
		TestFormValues(t, r, values{})
		_, _ = fmt.Fprint(w, string(jsonBytes))
	})

	c := GithubClientImpl{
		Client: client,
		cfg: &conf.GlobalConfig{
			GitProviderConfig: conf.GitProviderConfig{
				OrgLevelWebhook: true,
				OrgName:         "test",
				WebhookURL:      "https://bla.com",
			},
		},
	}
	ctx := context.Background()

	//
	// Execute
	//
	hooks, isEnabled := isOrgWebhookEnabled(ctx, &c)

	//
	// Assert
	//
	assert := assertion.New(t)
	assert.True(isEnabled)
	assert.NotNil(t, hooks)
}

func TestIsRepoWebhookEnabled(t *testing.T) {
	//
	// Prepare
	//
	client, mux, _, teardown := setup()
	defer teardown()

	config := make(map[string]interface{})
	config["url"] = "https://bla.com"
	Hooks := github.Hook{
		Active: utils.BPtr(true),
		Name:   utils.SPtr("web"),
		Config: config,
	}
	jsonBytes, _ := json.Marshal(&[]github.Hook{Hooks})

	mux.HandleFunc("/repos/test/test-repo2/hooks", func(w http.ResponseWriter, r *http.Request) {
		TestMethod(t, r, "GET")
		TestFormValues(t, r, values{})
		_, _ = fmt.Fprint(w, string(jsonBytes))
	})

	c := GithubClientImpl{
		Client: client,
		cfg: &conf.GlobalConfig{
			GitProviderConfig: conf.GitProviderConfig{
				OrgLevelWebhook: false,
				OrgName:         "test",
				WebhookURL:      "https://bla.com",
				RepoList:        "test-repo1,test-repo2",
			},
		},
	}
	ctx := context.Background()

	//
	// Execute
	//
	hook, isEnabled := isRepoWebhookEnabled(ctx, &c, "test-repo2")

	//
	// Assert
	//
	assert := assertion.New(t)
	assert.True(isEnabled)
	assert.NotNil(t, hook)

	//
	// Execute
	//
	hook, isEnabled = isRepoWebhookEnabled(ctx, &c, "test-repo3")

	//
	// Assert
	//
	assert = assertion.New(t)
	assert.False(isEnabled)
	assert.NotNil(t, hook)
}
