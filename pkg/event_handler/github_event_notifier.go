package event_handler

import (
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
)

type githubNotifier struct {
	cfg     *conf.GlobalConfig
	clients *clients.Clients
}

func NewGithubEventNotifier(cfg *conf.GlobalConfig, clients *clients.Clients) EventNotifier {
	return &githubNotifier{
		cfg:     cfg,
		clients: clients,
	}
}

func (gn *githubNotifier) notify(workflow v1alpha1.Workflow) {
	fmt.Printf("Notifing workflow, %s", workflow.GetName())
}
