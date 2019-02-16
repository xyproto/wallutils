package gnometimed

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/xyproto/crossfade"
	"github.com/xyproto/event"
)

// EventLoop will start the event loop for this GNOME Timed Wallpaper
func (gw *Wallpaper) EventLoop(verbose bool, setWallpaperFunc func(string) error) error {

	if verbose {
		fmt.Println("Using the GNOME Timed Wallpaper format")
	}

	// Convert to a SimpleTimedWallpaper, only for setting the initial wallpaper
	stw, err := GnomeToSimple(gw)
	if err != nil {
		return err
	}
	stw.SetInitialWallpaper(verbose, setWallpaperFunc)

	// Start the event loop
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
			if verbose {
				fmt.Printf("Registering static event at %s for setting %s\n", c(eventTime), s.Filename)
			}

			// Place values into variables, before enclosing it in the function below.
			from := eventTime
			window := s.Duration()
			cooldown := window
			imageFilename := s.Filename

			// Register a static event
			eventloop.Add(event.New(from, window, cooldown, func() {
				if verbose {
					fmt.Printf("Triggered static wallpaper event at %s\n", c(from))
					fmt.Println("Window:", window)
					fmt.Println("Cooldown:", cooldown)
					fmt.Println("Filename:", imageFilename)
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
				if verbose {
					fmt.Printf("Setting %s.\n", imageFilename)
				}
				if err := setWallpaperFunc(imageFilename); err != nil {
					fmt.Fprintf(os.Stderr, "Could not set wallpaper: %v\n", err)
					return // return from anon func
				}
			}))

			// Increase the variable that keeps track of the time
			eventTime = eventTime.Add(window)

		} else if t, ok := eInterface.(GTransition); ok {
			if verbose {
				fmt.Printf("Registering transition at %s for transitioning from %s to %s.\n", c(eventTime), t.FromFilename, t.ToFilename)
			}

			// cross fade steps
			steps := 10

			from := eventTime
			window := t.Duration()
			cooldown := window / time.Duration(steps)
			upTo := eventTime.Add(window)
			tType := t.Type
			tFromFilename := t.FromFilename
			tToFilename := t.ToFilename
			loopWait := gw.LoopWait
			tempDir := ""
			var err error

			// Register a transition event
			eventloop.Add(event.New(from, window, cooldown, func() {
				progress := window - event.ToToday(upTo).Sub(event.ToToday(time.Now()))
				for progress > h24 {
					progress -= h24
				}
				for progress < 0 {
					progress += h24
				}
				ratio := float64(progress) / float64(window)
				if verbose {
					fmt.Printf("Triggered transition event at %s (%d%% complete)\n", c(from), int(ratio*100))
					fmt.Println("Progress:", progress)
					fmt.Println("Up to:", c(upTo))
					fmt.Println("Window:", window)
					fmt.Println("Cooldown:", cooldown)
					fmt.Println("Loop wait:", loopWait)
					fmt.Println("Transition type:", tType)
					fmt.Println("From filename", tFromFilename)
					fmt.Println("To filename", tToFilename)
				}

				if exists(tempDir) {
					// Clean up
					os.RemoveAll(tempDir)
				}

				tempDir, err = ioutil.TempDir("", "crossfade")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Could not create temporary directory: %v\n", err)
					return // return from anon func
				}
				// Prepare to write an image to the temporary directory
				tempImageFilename := filepath.Join(tempDir, "out.png") // .jpg is also possible

				if verbose {
					fmt.Println("Crossfading between images.")
				}

				// Crossfade and write the new image to the temporary directory
				if crossfade.Files(tFromFilename, tToFilename, tempImageFilename, ratio) != nil {
					fmt.Fprintf(os.Stderr, "Could not crossfade images in transition: %v\n", err)
					return // return from anon func
				}

				// Double check that the generated file exists
				if _, err := os.Stat(tempImageFilename); os.IsNotExist(err) {
					fmt.Fprintf(os.Stderr, "File does not exist: %s\n", tempImageFilename)
					return // return from anon func
				}

				// Set the desktop wallpaper, if possible
				if verbose {
					fmt.Printf("Setting %s.\n", tempImageFilename)
				}
				if err := setWallpaperFunc(tempImageFilename); err != nil {
					fmt.Fprintf(os.Stderr, "Could not set wallpaper: %v\n", err)
					return // return from anon func
				}

			}))

			// Increase the variable that keeps track of the time
			eventTime = eventTime.Add(window)
		} else {
			// This should never happen, it would be an implementation error
			panic("got an element that is not a GStatic and not a GTransition")
		}
	}

	// Endless loop! Will wait loopWait duration between each event loop cycle.
	eventloop.Go(gw.LoopWait)

	return nil
}
