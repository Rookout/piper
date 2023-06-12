package webhook_hanlder

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v3"

	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/git"
	"github.com/rookout/piper/pkg/utils"
)

type WebhookHandlerImpl struct {
	cfg      *conf.Config
	clients  *clients.Clients
	Triggers *[]Trigger
	Payload  *git.WebhookPayload
}

func NewWebhookHandler(cfg *conf.Config, clients *clients.Clients, payload *git.WebhookPayload) (*WebhookHandlerImpl, error) {
	var err error

	return &WebhookHandlerImpl{
		cfg:      cfg,
		clients:  clients,
		Triggers: &[]Trigger{},
		Payload:  payload,
	}, err
}

func (wh *WebhookHandlerImpl) RegisterTriggers() error {
	if !IsFileExists(wh, "", ".workflows") {
		return fmt.Errorf(".workflows folder does not exist in %s/%s", wh.Payload.Repo, wh.Payload.Branch)
	}

	if !IsFileExists(wh, ".workflows", "triggers.yaml") {
		return fmt.Errorf(".workflows/triggers.yaml file does not exist in %s/%s", wh.Payload.Repo, wh.Payload.Branch)
	}

	triggers, err := wh.clients.Git.GetFile(wh.Payload.Repo, wh.Payload.Branch, ".workflows/triggers.yaml")
	if err != nil {
		return fmt.Errorf("failed to get triggers content: %v", err)
	}

	log.Printf("triggers content is: \n %s \n", *triggers.Content)

	err = yaml.Unmarshal([]byte(*triggers.Content), wh.Triggers)
	if err != nil {
		return fmt.Errorf("failed to unmarshal triggers content: %v", err)
	}
	return nil
}

func (wh *WebhookHandlerImpl) ExecuteMatchingTriggers() error {
	triggered := false
	for _, trigger := range *wh.Triggers {
		if utils.IsElementMatch(wh.Payload.Branch, *trigger.Branches) && utils.IsElementMatch(wh.Payload.Event, *trigger.Events) {
			log.Printf("Trigger %s for branch %s triggered", wh.Payload.Event, wh.Payload.Branch)
			triggered = true
		}
	}
	if !triggered {
		return fmt.Errorf("no matching trigger found for %s in branch %s", wh.Payload.Event, wh.Payload.Branch)
	}
	return nil
}

func IsFileExists(wh *WebhookHandlerImpl, path string, file string) bool {
	files, err := wh.clients.Git.ListFiles(wh.Payload.Repo, wh.Payload.Branch, path)
	if err != nil {
		log.Printf("Error listing files in repo: %s branch: %s. %v", wh.Payload.Repo, wh.Payload.Branch, err)
		return false
	}
	if len(files) == 0 {
		log.Printf("Empty list of files in repo: %s branch: %s", wh.Payload.Repo, wh.Payload.Branch)
		return false
	}

	if utils.IsElementExists(files, file) {
		return true
	}

	return false
}
