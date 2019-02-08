package event

import (
	"time"
)

// EventLoop is a collection of events
type EventLoop []*Event

// NewEventLoop creates an empty collection of events
func NewEventLoop() *EventLoop {
	return &EventLoop{}
}

// Add adds an event to the collection
func (el *EventLoop) Add(e *Event) {
	*el = append(*el, e)
}

// Once runs a given action only once, within a 1 second window of time
func (el *EventLoop) Once(when time.Time, action func()) {
	// The window is also used as the cooldown
	el.Add(New(when, 1*time.Second, 1*time.Second, action))
}

// OnceWindow runs a given action only once, within a custom duration
func (el *EventLoop) OnceWindow(when time.Time, window time.Duration, action func()) {
	// The window is also used as the cooldown
	el.Add(New(when, window, window, action))
}

// Loop launches an event loop that will sleep the given duration at every loop.
func (el *EventLoop) Go(sleep time.Duration) {
	// Use an endless event loop
	for {
		// For each possible event
		for _, e := range *el {
			// Check if the event should trigger
			if e.ShouldTrigger() {
				// When triggering an event, run it in the background
				go e.Trigger()
			}
			// If the window for an event is in the past, remove it,
			// unless it is a clock-event.
		}
		// How long to sleep before checking again
		time.Sleep(sleep)
	}
}
