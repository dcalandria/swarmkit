package events

import "sync"

// Channel provides a sink that can be listened on. The writer and channel
// listener must operate in separate goroutines.
//
// Consumers should listen on Channel.C until Closed is closed.
type Channel struct {
	C chan Event

	closed chan struct{}
	once   sync.Once
}

// NewChannel returns a channel. If buffer is zero, the channel is
// unbuffered.
func NewChannel(buffer int) *Channel {
	return &Channel{
		C:      make(chan Event, buffer),
		closed: make(chan struct{}),
	}
}

// Done returns a channel that will always proceed once the sink is closed.
func (ch *Channel) Done() chan struct{} {
	return ch.closed
}

// Write the event to the channel. Must be called in a separate goroutine from
// the listener.
func (ch *Channel) Write(event Event) error {
	select {
	case ch.C <- event:
		return nil
	case <-ch.closed:
		return ErrSinkClosed
	}
}

// Close the channel sink.
func (ch *Channel) Close() error {
	ch.once.Do(func() {
		close(ch.closed)
	})

	return nil
}
