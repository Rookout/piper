package git

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v52/github"
	"github.com/rookout/piper/pkg/conf"
	assertion "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestIsOrgWebhookEnabled(t *testing.T) {
	//
	// Prepare
	//
	client, mux, _, teardown := setup()
	defer teardown()

	active := true
	hookName := "web"
	config := make(map[string]interface{})
	config["url"] = "https://bla.com"
	Hooks := github.Hook{
		Active: &active,
		Name:   &hookName,
		Config: config,
	}
	jsonBytes, _ := json.Marshal(&[]github.Hook{Hooks})

	mux.HandleFunc("/orgs/test/hooks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		_, _ = fmt.Fprint(w, string(jsonBytes))
	})

	c := GithubClientImpl{
		client: client,
		cfg: &conf.Config{
			GitConfig: conf.GitConfig{
				OrgName:    "test",
				WebhookURL: "https://bla.com",
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
