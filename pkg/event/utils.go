package event

import "time"

// ToToday moves the date of a given time.Time to today's date.
// The hour/minute/second is kept as it is.
func ToToday(d time.Time) time.Time {
	// Get hour, minute and second from the event
	hour, min, sec := d.Clock()

	// Get the current time and date
	now := time.Now()

	// Return a new time.Time
	return time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, now.Nanosecond(), now.Location())
}
