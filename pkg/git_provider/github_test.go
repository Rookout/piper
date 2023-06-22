package git_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-github/v52/github"
	"github.com/rookout/piper/pkg/conf"
	assertion "github.com/stretchr/testify/assert"
)

func TestListFiles(t *testing.T) {
	//
	// Prepare
	//
	client, mux, _, teardown := setup()
	defer teardown()

	contentName := ".workflows"
	contentType := "dir"
	contentPath := ""
	jsonBytes, _ := json.Marshal(&github.RepositoryContent{
		Type: &contentType,
		Name: &contentName,
		Path: &contentPath,
	})

	mux.HandleFunc("/repos/test/test-repo1/contents/?ref=branch1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
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

	//
	// Execute
	//
	actualContent, err := c.ListFiles(&ctx, "test-repo1", "branch1", "")
	expectedContent := `[{"name":".workflows"',"path":".workflows","type":"dir"}]`

	//
	// Assert
	//
	assert := assertion.New(t)
	assert.Equal(expectedContent, actualContent)
	assert.NotNil(t, err)
}
