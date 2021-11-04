package event

// Event is an interface that will trigger at a specific hour and minute.
// Trigger() is the function that will be triggered.
// JustOnce() is if it should only be triggered once, or each day.
type Event interface {
	Trigger() error
	Hour() int
	Minute() int
	JustOnce() bool
}
