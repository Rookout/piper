package listener

type EventPubSubExample struct {
	callbacks map[string][]func(eventData any)
}

func (e *EventPubSubExample) Subscribe(eventName string, callback func(eventData any)) error {
	e.callbacks[eventName] = append(e.callbacks[eventName], callback)
	return nil
}

func (e *EventPubSubExample) Publish(eventName string, eventData any) error {
	for _, callback := range e.callbacks[eventName] {
		callback(eventData)
	}

	return nil
}

func NewEventPubSubExample() *EventPubSubExample {
	return &EventPubSubExample{
		callbacks: make(map[string][]func(eventData any)),
	}
}
