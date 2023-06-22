package clients

import (
	"github.com/rookout/piper/pkg/git_provider"
	"github.com/rookout/piper/pkg/workflow_handler"
)

type Clients struct {
	GitProvider git_provider.Client
	Workflows   workflow_handler.WorkflowsClient
}
