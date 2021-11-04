# Event II

[![Build Status](https://travis-ci.com/xyproto/event2.svg?branch=master)](https://travis-ci.com/xyproto/event2) [![GoDoc](https://godoc.org/github.com/xyproto/event2?status.svg)](https://godoc.org/github.com/xyproto/event2) [![License](https://img.shields.io/badge/license-MIT-green.svg?style=flat)](https://raw.githubusercontent.com/xyproto/event2/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/event2)](https://goreportcard.com/report/github.com/xyproto/event2)

A simple event system, for triggering events at certain times. This is the successor of [event](https://github.com/xyproto/event), which was needlessly complex.

## Example use

**Leet o'clock**

```go
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/xyproto/event2"
)

func main() {
	// Create a new event system, with a loop iteration delay of 1 second
	eventSys := event2.NewSystem(1 * time.Second)
	// Add an event that will trigger every day at 13:37
	eventSys.ClockEvent(13, 37, func() error {
		fmt.Println("It's leet o'clock")
		return nil
	})
	// Run the event system (not verbose)
	eventSys.Run(false)
	// Wait endlessly
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
```

**Clock**

```go
package main

import (
	"fmt"
	"time"

	"github.com/xyproto/event2"
)

func clockSystem() *event2.EventSys {
	sys := event2.NewSystem(1 * time.Second)
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
	clockSystem().Run(false)
	// Wait endlessly while saying "tick" and "tock" every second
	for {
		fmt.Println("tick")
		time.Sleep(1 * time.Second)
		fmt.Println("tock")
		time.Sleep(1 * time.Second)
	}
}
```

## General info

* Version: 0.0.1
* License: MIT
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
