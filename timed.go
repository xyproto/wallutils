package monitor

import (
	"fmt"
	"github.com/xyproto/event"
	"os"
	"path/filepath"
	"time"
)

// clockDuration finds the length of time between two given times,
// ignoring year/month day. A positive duration number will be returned.
func clockDuration(a, b time.Time) time.Duration {
	now := time.Now()
	at := time.Date(now.Year(), now.Month(), now.Day(), a.Hour(), a.Minute(), a.Second(), a.Nanosecond(), a.Location())
	bt := time.Date(now.Year(), now.Month(), now.Day(), b.Hour(), b.Minute(), b.Second(), b.Nanosecond(), b.Location())
	if at.Sub(bt) > 0 {
		return at.Sub(bt)
	}
	return bt.Sub(at)
}

// c formats a timestamp as HH:MM
func c(t time.Time) string {
	return fmt.Sprintf("%.2d:%.2d", t.Hour(), t.Minute())
}

// SetTimedWallpaper runs an infinite event loop and changes wallpapers at specified times
func SetTimedWallpaper(gw *GnomeWallpaper, verbose bool) error {

	loopWait := 5 * time.Second

	if verbose {
		fmt.Println("Using timed wallpaper:", gw.CollectionName)
	}

	eventloop := event.NewLoop()

	// Get the start time for the wallpaper collection (which is offset by X
	// seconds per static wallpaper)
	startTime := gw.StartTime()

	// The start time of the timed wallpaper as a whole
	if verbose {
		fmt.Println("Timed wallpaper start time:", c(startTime))
	}

	totalElements := len(gw.Config.Statics) + len(gw.Config.Transitions)

	// Keep track of the total time. It is increased every time a new element duration is encountered.
	eventTime := startTime

	for i := 0; i < totalElements; i++ {
		// The duration of the event is specified in the XML file, but not when it should start

		// Get an element, by index. This is an interface{} and is expected to be a GStatic or a GTransition
		eInterface, err := gw.Config.Get(i)
		if err != nil {
			return err
		}
		if s, ok := eInterface.(GStatic); ok {
			window := s.Duration()
			// We have a static GNOME wallpaper element, with a duration and an image filename

			if verbose {
				fmt.Printf("Registering STATIC at %s for changing to %s\n", c(eventTime), s.Filename)
			}

			// Place values into variables, before enclosing it in the function below.
			from := eventTime
			cooldown := window
			//imageFilename := s.Filename
			eventloop.Add(event.New(from, window, cooldown, func() {

				// Enclose variable
				imageFilename := s.Filename

				if verbose {
					fmt.Println("TRIGGERED STATIC WALLPAPER EVENT")
					fmt.Println("FROM", c(from))
					fmt.Println("WINDOW", window)
					fmt.Println("COOLDOWN", cooldown)
					fmt.Println("FILENAME", imageFilename)
				}

				// Find the absolute path
				absImageFilename, err := filepath.Abs(imageFilename)
				if err == nil {
					imageFilename = absImageFilename
				}

				// Check that the file exists
				if _, err := os.Stat(imageFilename); os.IsNotExist(err) {
					fmt.Fprintf(os.Stderr, "File does not exist: %s\n", imageFilename)
					return // return from anon func
				}

				// Set the desktop wallpaper, if possible
				if err := SetWallpaper(imageFilename); err != nil {
					fmt.Fprintf(os.Stderr, "Could not set wallpaper: %s\n", err)
					return // return from anon func
				}
			}))
			// Increase the variable that keeps track of the time
			//eventTime = event.ToToday(eventTime.Add(window))
			eventTime = eventTime.Add(window)

		} else if t, ok := eInterface.(GTransition); ok {
			// Increase the variable that keeps track of the time
			window := t.Duration()

			// We have a GNOME wallpaper transition, with a duration, a type,
			// and two image filenames.

			if verbose {
				fmt.Printf("Registering TRANSITION at %s for transitioning to %s.\n", c(eventTime), t.ToFilename)
			}

			from := eventTime
			cooldown := window
			upTo := eventTime.Add(window)
			eventloop.Add(event.New(from, window, cooldown, event.ProgressWrapperInterval(from, upTo, loopWait, func(p float64) {

				// Enclose variables
				tType := t.Type
				tFromFilename := t.FromFilename
				tToFilename := t.ToFilename

				if verbose {
					fmt.Println("TRIGGERED TRANSITION EVENT")
					fmt.Println("TO IMPLEMENT: A smooth transition")
					fmt.Println("TYPE         ", tType)
					fmt.Println("FROM FILENAME", tFromFilename)
					fmt.Println("TO FILENAME  ", tToFilename)
					fmt.Printf("PERCENTAGE COMPLETE: %d%%\n", int(p*100))
					fmt.Println("FROM", c(from))
					fmt.Println("WINDOW", window)
					fmt.Println("COOLDOWN", cooldown)
					fmt.Println("EVENT TIME", c(eventTime))
					fmt.Println("UP TO", c(upTo))
					fmt.Println("LOOP WAIT", loopWait)
				}

				// TODO: Create a temporary image that is a mix between t.FromFilename and t.ToFilename, using p as the ratio
				imageFilename := tToFilename

				// Find the absolute path
				absImageFilename, err := filepath.Abs(imageFilename)
				if err == nil {
					imageFilename = absImageFilename
				}

				// Check that the file exists
				if _, err := os.Stat(imageFilename); os.IsNotExist(err) {
					fmt.Fprintf(os.Stderr, "File does not exist: %s\n", imageFilename)
					return // return from anon func
				}

				// Set the desktop wallpaper, if possible
				if err := SetWallpaper(imageFilename); err != nil {
					fmt.Fprintf(os.Stderr, "Could not set wallpaper: %s\n", err)
					return // return from anon func
				}
			})))

			//eventTime = event.ToToday(eventTime.Add(window))
			eventTime = eventTime.Add(window)
		} else {
			// This should never happen, it would be an implementation error
			panic("got an element that is not a GStatic and not a GTransition")
		}
	}

	// Endless loop! Will wait loopWait duration between each event loop cycle.
	eventloop.Go(loopWait)

	return nil
}
