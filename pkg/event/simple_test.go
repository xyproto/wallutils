package event

import (
	"testing"
	"time"
)

func TestNewSimpleEvent(t *testing.T) {
	sys := NewSystem(1 * time.Second)
	e := NewTestEvent()
	sys.Register(e)
	verbose := true
	sys.Run(verbose)
}
