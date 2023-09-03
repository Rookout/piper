package git_provider

import (
	"encoding/json"
	"fmt"
	"github.com/ktrysmt/go-bitbucket"
	"github.com/rookout/piper/pkg/conf"
	assertion "github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"net/http"
	"testing"
)

func TestBitbucketListFiles(t *testing.T) {
	// Prepare
	client, mux, _, teardown := setupBitbucket()
	defer teardown()

	repoContent := &bitbucket.RepositoryFile{
		Type: "file",
		Path: ".workflows/exit.yaml",
	}

	repoContent2 := &bitbucket.RepositoryFile{
		Type: "file",
		Path: ".workflows/main.yaml",
	}

	data := map[string]interface{}{"values": []bitbucket.RepositoryFile{*repoContent, *repoContent2}}
	jsonBytes, _ := json.Marshal(data)

	mux.HandleFunc("/repositories/test/test-repo1/src/branch1/.workflows/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		//testFormValues(t, r, values{})

		_, _ = fmt.Fprint(w, string(jsonBytes))
	})

	c := BitbucketClientImpl{
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
