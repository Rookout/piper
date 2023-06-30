package notifier

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExampleNotifier(t *testing.T) {

	exampleNotifier := NewChannelNotifier()

	exitSignal := make(chan bool)

	puller := func() {
		exampleNotifier.PullMessagesUntilSuccess()
		exitSignal <- true
	}

	worker := func(notifier Notifier) {
		_ = notifier.Progress("10% done")
		_ = notifier.Progress("50% done")
		_ = notifier.Failure("failed. retrying ...")
		_ = notifier.Success("task completed")
	}

	go puller()
	go worker(&exampleNotifier)

	<-exitSignal

	history := exampleNotifier.GetNotificationsHistory()

	assert.Equal(
		t,
		[]string{
			"10% done",
			"50% done",
			"failed. retrying ...",
			"task completed",
		},
		history,
	)
}
