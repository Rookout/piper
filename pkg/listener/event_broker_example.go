package listener

type EventBrokerExample struct {
	callbacks map[string][]func(eventData any)
}

func (e *EventBrokerExample) Subscribe(eventName string, callback func(eventData any)) error {
	e.callbacks[eventName] = append(e.callbacks[eventName], callback)
	return nil
}

func (e *EventBrokerExample) Publish(eventName string, eventData any) error {
	for _, callback := range e.callbacks[eventName] {
		callback(eventData)
	}

	return nil
}

func NewEventBrokerExample() *EventBrokerExample {
	return &EventBrokerExample{
		callbacks: make(map[string][]func(eventData any)),
	}
}
