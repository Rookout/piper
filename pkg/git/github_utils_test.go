package git

import (
	"context"
	"fmt"
	"github.com/rookout/piper/pkg/conf"
	"net/http"
	"testing"
)

func TestIsOrgWebhookEnabled(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/test/hooks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	c := GithubClientImpl{
		client: client,
		cfg: &conf.Config{
			GitConfig: conf.GitConfig{
				OrgName: "test",
			},
		},
	}

	ctx := context.Background()
	isOrgWebhookEnabled(ctx, &c)

}
