package git_provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ktrysmt/go-bitbucket"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/utils"
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

func TestBitbucketSetStatus(t *testing.T) {
	// Prepare
	ctx := context.Background()
	assert := assertion.New(t)
	client, mux, _, teardown := setupBitbucket()
	defer teardown()

	mux.HandleFunc("/repositories/test/test-repo1/commit/test-commit/statuses/build", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{})

		w.WriteHeader(http.StatusCreated)
		jsonBytes := []byte(`{"status": "ok"}`)
		_, _ = fmt.Fprint(w, string(jsonBytes))
	})

	c := BitbucketClientImpl{
		client: client,
		cfg: &conf.GlobalConfig{
			GitProviderConfig: conf.GitProviderConfig{
				Provider:        "bitbucket",
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
