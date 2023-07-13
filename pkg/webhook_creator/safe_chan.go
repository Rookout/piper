package webhook_creator

import "sync"

type SafeChannel struct {
	C      chan struct{}
	closed bool
	mutex  sync.Mutex
}

func NewSafeChannel() *SafeChannel {
	return &SafeChannel{C: make(chan struct{})}
}

func (mc *SafeChannel) SafeClose() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	if !mc.closed {
		close(mc.C)
		mc.closed = true
	}
}

func (mc *SafeChannel) IsClosed() bool {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	return mc.closed
}
