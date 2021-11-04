package event

import (
	"fmt"
	"testing"
	"time"
)

func TestEveryMinute(t *testing.T) {
	sys := NewSystem(1 * time.Second)
	now := time.Now()
	i := 0
	n := 1 // Could be any number of events to trigger, with a minute between
	sys.EveryMinute(now.Hour(), now.Minute(), n, func() error {
		fmt.Printf("WAIT A MINUTE #%d\n", i+1)
		i++
		return nil
	})
	verbose := true
	sys.Run(verbose)
	time.Sleep(1 * time.Second)
}
