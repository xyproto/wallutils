package event

import (
	"time"
)

// ProgressWrapper wraps a function that takes a float64 in such a way that
// a float between 0 and 1 is passed in when it is called. The float signifies
// how far time has progressed from the given `from`  time up to the given
// `upto` time. The wrapped function that takes no arguments is returned.
// The final "100%" (1.0) value will never be sent to the wrapped function,
// because the `upto` time is exclusive, not inclusive.
func ProgressWrapper(from, upto time.Time, progressFunction func(float64)) func() {
	start := time.Now()
	duration := upto.Sub(from)
	// Wrap the given function in a function that can measure the rate of progress
	return func() {
		start := start
		duration := duration
		passed := time.Now().Sub(start)
		ratio := 0.0
		if duration > 0 {
			ratio = float64(passed) / float64(duration)
		}
		// Clamp the ratio
		if ratio > 1.0 {
			ratio = 1.0
		}
		// Call the wrapped function, with an appropriate ratio
		progressFunction(ratio)
	}
}

// ProgressWrapperInterval behaves like ProgressWrapper, except that an
// interval is given, that lets the function receive a "1.0" value at the
// final run, when progress is complete. The interval duration is subtracted
// from the `upto` time, at the time when the progress float is calculated.
func ProgressWrapperInterval(from, upto time.Time, interval time.Duration, progressFunction func(float64)) func() {
	start := time.Now()
	duration := upto.Sub(from) - interval
	// Wrap the given function in a function that can measure the rate of progress
	return func() {
		start := start
		duration := duration
		passed := time.Now().Sub(start)
		ratio := 0.0
		if duration > 0 {
			ratio = float64(passed) / float64(duration)
		}
		// Clamp the ratio
		if ratio > 1.0 {
			ratio = 1.0
		}
		// Call the wrapped function, with an appropriate ratio
		progressFunction(ratio)
	}
}
