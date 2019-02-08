package event

import (
	"time"
)

// Events is a collection of events
type Events []*Event

// NewEvents creates an empty collection of events
func NewEvents() *Events {
	return &Events{}
}

// Add adds an event to the collection
func (es *Events) Add(e *Event) {
	*es = append(*es, e)
}

// Once runs a given action only once, within a 1 second window of time
func (es *Events) Once(when time.Time, action func()) {
	// The window is also used as the cooldown
	es.Add(New(when, 1*time.Second, 1*time.Second, action))
}

// OnceWindow runs a given action only once, within a custom duration
func (es *Events) OnceWindow(when time.Time, window time.Duration, action func()) {
	// The window is also used as the cooldown
	es.Add(New(when, window, window, action))
}

// Loop launches an event loop that will sleep the given duration at every loop.
func (es *Events) Loop(sleep time.Duration) {
	// Use an endless event loop
	for {
		// For each possible event
		for _, e := range *es {
			// Check if the event should trigger
			if e.ShouldTrigger() {
				// When triggering an event, run it in the background
				go e.Trigger()
			}
		}
		// How long to sleep before checking again
		time.Sleep(sleep)
	}
}
