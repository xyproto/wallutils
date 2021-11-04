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
