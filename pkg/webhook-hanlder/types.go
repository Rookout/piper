package webhook_hanlder

type Trigger struct {
	events   []string `json:"events"`
	branches []string `json:"branches"`
	onStart  []string `json:"execute"`
	onExit   []string `json:"on_exit"`
}

type WebhookHandler interface {
	RegisterTriggers() error
	ExecuteMatchingTriggers() error
}
