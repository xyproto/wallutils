package monitor

import (
	"fmt"
	"github.com/xyproto/crossfade"
	"github.com/xyproto/event"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// SetSimpleTimedWallpaper runs an infinite event loop and changes wallpapers at specified times
func SetSimpleTimedWallpaper(stw *SimpleTimedWallpaper, verbose bool) error {

	loopWait := 5 * time.Second

	if verbose {
		fmt.Println("Using the Simple Timed Wallpaper format.")
	}

	eventloop := event.NewLoop()

	for _, s := range stw.Statics {
		if verbose {
			fmt.Printf("Registering STATIC event at %s for setting %s\n", c(s.At), s.Filename)
		}

		// Place values into variables, before enclosing it in the function below.
		from := s.At
		window := stw.UntilNext(s.At) // duration until next event start
		cooldown := window
		imageFilename := s.Filename

		// Register a static event
		eventloop.Add(event.New(from, window, cooldown, func() {
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
	}
	for _, t := range stw.Transitions {
		if verbose {
			fmt.Printf("Registering TRANSITION at %s for transitioning to %s.\n", c(t.From), t.ToFilename)
		}

		// cross fade steps
		steps := 10

		// Set variables
		from := t.From
		window := t.Duration()
		cooldown := window / time.Duration(steps)
		upTo := from.Add(window)
		tType := t.Type
		tFromFilename := t.FromFilename
		tToFilename := t.ToFilename

		// Register a transition event
		eventloop.Add(event.New(from, window, cooldown, event.ProgressWrapperInterval(from, upTo, loopWait, func(ratio float64) {
			if verbose {
				fmt.Println("TRIGGERED TRANSITION EVENT")
				fmt.Println("TYPE         ", tType)
				fmt.Println("FROM FILENAME", tFromFilename)
				fmt.Println("TO FILENAME  ", tToFilename)
				fmt.Printf("PERCENTAGE COMPLETE: %d%%\n", int(ratio*100))
				fmt.Println("FROM", c(from))
				fmt.Println("UP TO", c(upTo))
				fmt.Println("WINDOW", window)
				fmt.Println("COOLDOWN", cooldown)
				fmt.Println("LOOP WAIT", loopWait)
			}

			tempDir, err := ioutil.TempDir("", "crossfade")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not create temporary directory: %s\n", err)
				return // return from anon func
			}
			// TODO: Find out if it is safe to remove the wallpaper image while it is in use, or not
			defer os.RemoveAll(tempDir) // clean up

			// Prepare to write an image to the temporary directory
			tempImageFilename := filepath.Join(tempDir, "out.png") // .jpg is also possible

			// Crossfade and write the new image to the temporary directory
			if crossfade.Files(tFromFilename, tToFilename, tempImageFilename, ratio) != nil {
				fmt.Fprintf(os.Stderr, "Could not crossfade images in transition: %s\n", err)
				return // return from anon func
			}

			// Double check that the generated file exists
			if _, err := os.Stat(tempImageFilename); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "File does not exist: %s\n", tempImageFilename)
				return // return from anon func
			}

			// Set the desktop wallpaper, if possible
			if err := SetWallpaper(tempImageFilename); err != nil {
				fmt.Fprintf(os.Stderr, "Could not set wallpaper: %s\n", err)
				return // return from anon func
			}
		})))
	}

	// Endless loop! Will wait loopWait duration between each event loop cycle.
	eventloop.Go(loopWait)
	return nil
}

// SetGnomeTimedWallpaper runs an infinite event loop and changes wallpapers at specified times
func SetGnomeTimedWallpaper(gw *GnomeWallpaper, verbose bool) error {

	loopWait := 5 * time.Second

	if verbose {
		fmt.Println("Using the GNOME Timed Wallpaper format")
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
			if verbose {
				fmt.Printf("Registering STATIC at %s for setting %s\n", c(eventTime), s.Filename)
			}

			// Place values into variables, before enclosing it in the function below.
			from := eventTime
			window := s.Duration()
			cooldown := window
			imageFilename := s.Filename

			// Register a static event
			eventloop.Add(event.New(from, window, cooldown, func() {
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
			eventTime = eventTime.Add(window)

		} else if t, ok := eInterface.(GTransition); ok {
			if verbose {
				fmt.Printf("Registering TRANSITION at %s for transitioning to %s.\n", c(eventTime), t.ToFilename)
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

			// Register a transition event
			eventloop.Add(event.New(from, window, cooldown, event.ProgressWrapperInterval(from, upTo, loopWait, func(ratio float64) {
				if verbose {
					fmt.Println("TRIGGERED TRANSITION EVENT")
					fmt.Println("TYPE         ", tType)
					fmt.Println("FROM FILENAME", tFromFilename)
					fmt.Println("TO FILENAME  ", tToFilename)
					fmt.Printf("PERCENTAGE COMPLETE: %d%%\n", int(ratio*100))
					fmt.Println("FROM", c(from))
					fmt.Println("WINDOW", window)
					fmt.Println("COOLDOWN", cooldown)
					fmt.Println("EVENT TIME", c(eventTime))
					fmt.Println("UP TO", c(upTo))
					fmt.Println("LOOP WAIT", loopWait)
				}

				tempDir, err := ioutil.TempDir("", "crossfade")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Could not create temporary directory: %s\n", err)
					return // return from anon func
				}
				// TODO: Find out if it is safe to remove the wallpaper image while it is in use, or not
				defer os.RemoveAll(tempDir) // clean up

				// Prepare to write an image to the temporary directory
				tempImageFilename := filepath.Join(tempDir, "out.png") // .jpg is also possible

				// Crossfade and write the new image to the temporary directory
				if crossfade.Files(tFromFilename, tToFilename, tempImageFilename, ratio) != nil {
					fmt.Fprintf(os.Stderr, "Could not crossfade images in transition: %s\n", err)
					return // return from anon func
				}

				// Double check that the generated file exists
				if _, err := os.Stat(tempImageFilename); os.IsNotExist(err) {
					fmt.Fprintf(os.Stderr, "File does not exist: %s\n", tempImageFilename)
					return // return from anon func
				}

				// Set the desktop wallpaper, if possible
				if err := SetWallpaper(tempImageFilename); err != nil {
					fmt.Fprintf(os.Stderr, "Could not set wallpaper: %s\n", err)
					return // return from anon func
				}
			})))

			// Increase the variable that keeps track of the time
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
