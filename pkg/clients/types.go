package clients

import (
	"github.com/rookout/piper/pkg/git"
	workflowHandler "github.com/rookout/piper/pkg/workflow-handler"
)

type Clients struct {
	Git       git.Client
	Workflows workflowHandler.WorkflowsClient
}
