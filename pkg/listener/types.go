package listener

type EventBroker interface {
	Subscribe(eventName string, callback func(eventData any)) error
	Publish(eventName string, eventData any) error
}
