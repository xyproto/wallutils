package event

import (
	"log"
	"sync"
	"time"
)

// EventSys represents an event system.
// * granularity is how long it should wait at each loop in the main loop
// * events is a map from Event to a margin of error on each side of the time
//   when the event should kick in
// * coolOffGranularity is how long the system should wait per cool-off
//   loop iteration
type EventSys struct {
	events             []Event
	granularity        time.Duration
	coolOffGranularity time.Duration
}

// CoolOff is a list of all events that should not be triggered just yet
var (
	coolOff    []Event
	coolOffMut sync.Mutex
)

func (sys *EventSys) coolOffLoop() {
	// Every N seconds, remove the first entry from the coolOff slice
	for {
		coolOffMut.Lock()
		if len(coolOff) > 0 {
			// If the event should be ran just once, move it to the back of the queue
			if coolOff[0].JustOnce() {
				tmp := coolOff[0]
				coolOff = coolOff[1:]
				coolOff = append(coolOff, tmp)
			} else {
				coolOff = coolOff[1:]
			}
		}
		coolOffMut.Unlock()
		time.Sleep(sys.coolOffGranularity)
	}
}

// NewSystem creates a new event system, where events can be registered
// and the event loop can be run. loopSleep is how long the event loop
// should sleep at every iteration.
func NewSystem(loopSleep time.Duration) *EventSys {
	events := make([]Event, 0)
	granularity := loopSleep
	coolOffDuration := time.Minute * 5
	return &EventSys{events, granularity, coolOffDuration}
}

// Register will register an event with the event system.
func (sys *EventSys) Register(event Event) {
	// TODO: Not thread safe? Add a mutex?
	sys.events = append(sys.events, event)
}

// eventLoop will run the event system endlessly, in the foreground
func (sys *EventSys) eventLoop(verbose bool) error {
	for {
		// Check if any events should kick in at this point in time +- error margin, in seconds
		now := time.Now()
	NEXT_EVENT:
		for _, event := range sys.events {
			// If the event is in the coolOff slice, skip for now
			for _, coolOffEvent := range coolOff {
				if coolOffEvent == event {
					//if verbose {
					//log.Println("Skipping event that is in the cool-off period")
					//}
					continue NEXT_EVENT
				}
			}
			if now.Hour() == event.Hour() && now.Minute() == event.Minute() {
				if verbose {
					log.Printf("Trigger event at %02d:%02d\n", now.Hour(), now.Minute())
				}
				if event.Trigger() != nil {
					// TODO: Do something sensible if the trigger fails?
					if verbose {
						log.Println("Event failed")
					}
				}
				// Placing in the CoolOff slice,
				// which is handled by the cooloff-system
				coolOffMut.Lock()
				coolOff = append(coolOff, event)
				coolOffMut.Unlock()
			}
		}
		time.Sleep(sys.granularity)
	}
}

// RunBackground will start the event system in the background and immediately return
func (sys *EventSys) RunBackground(verbose bool) {
	go sys.coolOffLoop()
	go sys.eventLoop(verbose)
}

// Run will start the event system in the foreground and never return
func (sys *EventSys) Run(verbose bool) {
	sys.RunBackground(verbose)
	// Wait endlessly
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

// Reset will remove a given event from the cool-off queue
// TODO: This function needs testing!
func Reset(event Event) {
	coolOffMut.Lock()
	s := coolOff
	var i int
	for i2, e := range coolOff {
		if e == event {
			i = i2
			break
		}
	}
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	coolOff = s[:len(s)-1]
	coolOffMut.Unlock()
}

// SimpleEvent creates and registers an event that should happen in a
// certain amount of time from now, then may optionally be repeated at every
// matching hour and minute every 24 hours, if "once" is false.
func (sys *EventSys) SimpleEvent(in time.Duration, once bool, f func() error) {
	sys.Register(NewSimpleEvent(in, once, f))
}

// ClockEvent creates and registers an event that should happen at every HH:MM
func (sys *EventSys) ClockEvent(h, m int, f func() error) {
	sys.Register(NewClockEvent(h, m, f))
}

// EveryMinute will trigger an event every minute for n minutes, starting from h:m
func (sys *EventSys) EveryMinute(h, m, n int, f func() error) {
	for i := 0; i < n; i++ {
		sys.Register(NewClockEvent(h, m, f))
		m++
		if m >= 60 {
			h++
			m = 0
		}
		if h >= 24 {
			h = 0
		}
	}
}
