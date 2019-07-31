package event

import (
	"time"
)

// Loop is a collection of events
type Loop []*Event

// NewLoop creates an empty collection of events
func NewLoop() *Loop {
	return &Loop{}
}

// Add adds an event to the collection
func (l *Loop) Add(e *Event) {
	*l = append(*l, e)
}

// Once runs a given action only once, within a 1 second window of time
func (l *Loop) Once(when time.Time, action func()) {
	// The window is also used as the cooldown
	l.Add(NewOnce(when, 1*time.Second, 1*time.Second, action))
}

// OnceWindow runs a given action only once, within a custom duration
func (l *Loop) OnceWindow(when time.Time, window time.Duration, action func()) {
	// The window is also used as the cooldown
	cooldown := window
	l.Add(NewOnce(when, window, cooldown, action))
}

// Go launches an event loop that will sleep the given duration at every loop.
func (l *Loop) Go(sleep time.Duration) {
	// Use an endless event loop
	for {
		// For each possible event
		for _, e := range *l {
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
