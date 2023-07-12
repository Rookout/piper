package notifier

import "fmt"

type ChannelNotifier struct {
	successChan  chan string
	failureChan  chan string
	progressChan chan string
	history      []string
}

func NewChannelNotifier() ChannelNotifier {
	return ChannelNotifier{
		successChan:  make(chan string),
		failureChan:  make(chan string),
		progressChan: make(chan string),
		history:      make([]string, 0),
	}
}

func (c *ChannelNotifier) Success(msg string) error {
	c.successChan <- msg
	return nil
}

func (c *ChannelNotifier) Failure(msg string) error {
	c.failureChan <- msg
	return nil
}

func (c *ChannelNotifier) Progress(msg string) error {
	c.progressChan <- msg
	return nil
}

func (c *ChannelNotifier) GetNotificationsHistory() []string {
	return c.history
}

func (c *ChannelNotifier) PullMessagesUntilSuccess() {
	for {
		select {
		case msg := <-c.progressChan:
			fmt.Printf("progress notification: %s\n", msg)
			c.history = append(c.history, msg)

		case msg := <-c.successChan:
			fmt.Printf("success notification: %s\n", msg)
			c.history = append(c.history, msg)
			return

		case msg := <-c.failureChan:
			fmt.Printf("failure notification: %s\n", msg)
			c.history = append(c.history, msg)
		}
	}
}
