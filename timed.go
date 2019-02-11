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
		// TODO: The duration of a static is until the next event, not 15 minutes!
		window := 15 * time.Minute // s.Duration()
		eventTime := s.At

		// We have a static wallpaper element, with an image filename
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
		eventTime = eventTime.Add(window)

	}
	for _, t := range stw.Transitions {
		window := t.Duration()
		eventTime := t.From

		// We have a GNOME wallpaper transition, with a duration, a type,
		// and two image filenames.

		if verbose {
			fmt.Printf("Registering TRANSITION at %s for transitioning to %s.\n", c(eventTime), t.ToFilename)
		}

		// cross fade steps
		steps := 10

		from := eventTime
		cooldown := window / time.Duration(steps)
		upTo := eventTime.Add(window)
		eventloop.Add(event.New(from, window, cooldown, event.ProgressWrapperInterval(from, upTo, loopWait, func(p float64) {

			// Enclose variables
			tType := t.Type
			tFromFilename := t.FromFilename
			tToFilename := t.ToFilename
			ratio := p

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
			eventTime = eventTime.Add(window)

		} else if t, ok := eInterface.(GTransition); ok {
			window := t.Duration()

			// We have a GNOME wallpaper transition, with a duration, a type,
			// and two image filenames.

			if verbose {
				fmt.Printf("Registering TRANSITION at %s for transitioning to %s.\n", c(eventTime), t.ToFilename)
			}

			// cross fade steps
			steps := 10

			from := eventTime
			cooldown := window / time.Duration(steps)
			upTo := eventTime.Add(window)
			eventloop.Add(event.New(from, window, cooldown, event.ProgressWrapperInterval(from, upTo, loopWait, func(p float64) {

				// Enclose variables
				tType := t.Type
				tFromFilename := t.FromFilename
				tToFilename := t.ToFilename
				ratio := p

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
