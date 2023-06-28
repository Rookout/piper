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

func TestListFiles(t *testing.T) {
	// Prepare
	client, mux, _, teardown := Setup()
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
		TestMethod(t, r, "GET")
		//TestFormValues(t, r, values{})

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
		Client: client,
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
