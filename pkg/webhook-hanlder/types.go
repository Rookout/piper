package webhook_hanlder

type Trigger struct {
	events   []string `json:"events"`
	branches []string `json:"branches"`
	execute  []string `json:"execute"`
	onExit   []string `json:"on_exit"`
}

type WebhookHandler interface {
	RegisterTriggers(triggers *[]Trigger) error
	ExecuteMatchingTriggers(event string, branch string) error
}
