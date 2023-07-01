package listener

type SimplePubSub struct {
	callbacks map[string][]func(eventData any)
}

func (e *SimplePubSub) Subscribe(eventName string, callback func(eventData any)) error {
	e.callbacks[eventName] = append(e.callbacks[eventName], callback)
	return nil
}

func (e *SimplePubSub) Publish(eventName string, eventData any) error {
	for _, callback := range e.callbacks[eventName] {
		callback(eventData)
	}

	return nil
}

func NewSimplePubSub() *SimplePubSub {
	return &SimplePubSub{
		callbacks: make(map[string][]func(eventData any)),
	}
}
