package webhook_hanlder

type Trigger struct {
	Events   *[]string `yaml:"events"`
	Branches *[]string `yaml:"branches"`
	OnStart  *[]string `yaml:"onStart"`
	OnExit   *[]string `yaml:"onExit"`
}

type WebhookHandler interface {
	RegisterTriggers() error
	ExecuteMatchingTriggers() error
}
