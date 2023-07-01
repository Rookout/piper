package listener

type PubSub interface {
	Publish(eventName string, eventData any) error
	Subscriber
}

type Subscriber interface {
	Subscribe(eventName string, callback func(eventData any)) error
}
