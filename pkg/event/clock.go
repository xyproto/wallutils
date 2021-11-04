package event

// NewClockEvent will create a simple event that can be triggered
// The event will trigger every time the hour and minute matches the one from time.Now()
func NewClockEvent(h, m int, f func() error) *SimpleEvent {
	return &SimpleEvent{h, m, false, f}
}
