# Event

A simple event system, for triggering events at certain times.

## Example use

**Leet o'clock**

```go
package main

import (
	"fmt"
	"time"

	"github.com/xyproto/wallutils/pkg/event"
)

func main() {
	// Create a new event system, with a loop iteration delay of 1 second
	eventSys := event.NewSystem(1 * time.Second)
	// Add an event that will trigger every day at 13:37
	eventSys.ClockEvent(13, 37, func() error {
		fmt.Println("It's leet o'clock")
		return nil
	})
	// Run the event system (not verbose)
	eventSys.Run(false)
}
```

**Clock**

```go
package main

import (
	"fmt"
	"time"

	"github.com/xyproto/wallutils/pkg/event"
)

func clockSystem() *event.EventSys {
	sys := event.NewSystem(1 * time.Second)
	for hour := 0; hour < 24; hour++ {
		for minute := 0; minute < 60; minute++ {
			// Create new variables that can be closed over by the new function below
			hour := hour
			minute := minute
			// Create a new event that will trigger at the specified hour and minute
			sys.ClockEvent(hour, minute, func() error {
				fmt.Printf("The clock is %02d:%02d\n", hour, minute)
				return nil
			})
		}
	}
	return sys
}

func main() {
	// Start the event system that will trigger an event at every minute
	clockSystem().RunBackground(false)
	// Wait endlessly while saying "tick" and "tock" every second
	for {
		fmt.Println("tick")
		time.Sleep(1 * time.Second)
		fmt.Println("tock")
		time.Sleep(1 * time.Second)
	}
}
```
