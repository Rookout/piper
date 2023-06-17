package webhook_hanlder

import (
	"context"
	"fmt"
	"github.com/rookout/piper/pkg/clients"
	"github.com/rookout/piper/pkg/common"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/git"
	"github.com/rookout/piper/pkg/utils"
	"gopkg.in/yaml.v3"
	"log"
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

func (wh *WebhookHandlerImpl) RegisterTriggers(ctx *context.Context) error {
	if !IsFileExists(ctx, wh, "", ".workflows") {
		return fmt.Errorf(".workflows folder does not exist in %s/%s", wh.Payload.Repo, wh.Payload.Branch)
	}

	if !IsFileExists(ctx, wh, ".workflows", "triggers.yaml") {
		return fmt.Errorf(".workflows/triggers.yaml file does not exist in %s/%s", wh.Payload.Repo, wh.Payload.Branch)
	}

	triggers, err := wh.clients.Git.GetFile(ctx, wh.Payload.Repo, wh.Payload.Branch, ".workflows/triggers.yaml")
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

func (wh *WebhookHandlerImpl) PrepareBatchForMatchingTriggers(ctx *context.Context) ([]*common.WorkflowsBatch, error) {
	triggered := false
	var workflowBatches []*common.WorkflowsBatch
	for _, trigger := range *wh.Triggers {
		if trigger.Branches == nil {
			return nil, fmt.Errorf("trigger from repo %s branch %s missing branch field", wh.Payload.Repo, wh.Payload.Branch)
		}
		if trigger.Events == nil {
			return nil, fmt.Errorf("trigger from repo %s branch %s missing event field", wh.Payload.Repo, wh.Payload.Branch)
		}
		if utils.IsElementMatch(wh.Payload.Branch, *trigger.Branches) && utils.IsElementMatch(wh.Payload.Event, *trigger.Events) {
			log.Printf(
				"Triggering event %s for repo %s branch %s are triggered.",
				wh.Payload.Event,
				wh.Payload.Repo,
				wh.Payload.Branch,
			)
			triggered = true
			onStartFiles, err := wh.clients.Git.GetFiles(
				ctx,
				wh.Payload.Repo,
				wh.Payload.Branch,
				utils.AddPrefixToList(*trigger.OnStart, ".workflows/"),
			)
			if len(onStartFiles) == 0 {
				return nil, fmt.Errorf("one or more of onStart: %s files found", *trigger.OnStart)
			}
			if err != nil {
				return nil, err
			}

			onExitFiles := make([]*git.CommitFile, 0)
			if trigger.OnExit != nil {
				onExitFiles, err = wh.clients.Git.GetFiles(
					ctx,
					wh.Payload.Repo,
					wh.Payload.Branch,
					utils.AddPrefixToList(*trigger.OnExit, ".workflows/"),
				)
				if len(onExitFiles) == 0 {
					log.Printf("onExist: %s files not found", *trigger.OnExit)
				}
				if err != nil {
					return nil, err
				}
			}

			templatesFiles := make([]*git.CommitFile, 0)
			if trigger.Templates != nil {
				templatesFiles, err = wh.clients.Git.GetFiles(
					ctx,
					wh.Payload.Repo,
					wh.Payload.Branch,
					utils.AddPrefixToList(*trigger.Templates, ".workflows/"),
				)
				if len(templatesFiles) == 0 {
					log.Printf("parameters: %s files not found", *trigger.Templates)
				}
				if err != nil {
					return nil, err
				}
			}

			parameters, err := wh.clients.Git.GetFile(
				ctx,
				wh.Payload.Repo,
				wh.Payload.Branch,
				".workflows/parameters.yaml",
			)
			if err != nil {
				return nil, err
			}
			if parameters == nil {
				log.Printf("parameters.yaml not found in repo: %s branch %s", wh.Payload.Repo, wh.Payload.Branch)
			}

			workflowBatches = append(workflowBatches, &common.WorkflowsBatch{
				OnStart:    onStartFiles,
				OnExit:     onExitFiles,
				Templates:  templatesFiles,
				Parameters: parameters,
				Config:     &trigger.Config,
				Payload:    wh.Payload,
			})
		}
	}
	if !triggered {
		return nil, fmt.Errorf("no matching trigger found for %s in branch %s", wh.Payload.Event, wh.Payload.Branch)
	}
	return workflowBatches, nil
}

func IsFileExists(ctx *context.Context, wh *WebhookHandlerImpl, path string, file string) bool {
	files, err := wh.clients.Git.ListFiles(ctx, wh.Payload.Repo, wh.Payload.Branch, path)
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

func HandleWebhook(ctx *context.Context, wh *WebhookHandlerImpl) ([]*common.WorkflowsBatch, error) {
	err := wh.RegisterTriggers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to register triggers, error: %v", err)
	} else {
		log.Printf("successfully registered triggers for repo: %s branch: %s", wh.Payload.Repo, wh.Payload.Branch)
	}

	workflowsBatches, err := wh.PrepareBatchForMatchingTriggers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare matching triggers, error: %v", err)
	}

	if len(workflowsBatches) == 0 {
		log.Printf("no workflows to execute")
		return nil, fmt.Errorf("no workflows to execute for repo: %s branch: %s",
			wh.Payload.Repo,
			wh.Payload.Branch,
		)
	}
	return workflowsBatches, nil
}
