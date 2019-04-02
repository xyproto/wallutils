package gnometimed

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xyproto/crossfade"
	"github.com/xyproto/event"
)

// EventLoop will start the event loop for this GNOME Timed Wallpaper
func (gtw *Wallpaper) EventLoop(verbose bool, setWallpaperFunc func(string) error) error {

	if verbose {
		fmt.Println("Using the GNOME Timed Wallpaper format")
	}

	// Convert to a SimpleTimedWallpaper, only for setting the initial wallpaper
	stw, err := GnomeToSimple(gtw)
	if err != nil {
		return err
	}

	if err := stw.SetInitialWallpaper(verbose, setWallpaperFunc); err != nil {
		return err
	}

	// Start the event loop
	eventloop := event.NewLoop()

	// Get the start time for the wallpaper collection (which is offset by X
	// seconds per static wallpaper)
	startTime := gtw.StartTime()

	// The start time of the timed wallpaper as a whole
	if verbose {
		fmt.Println("Timed wallpaper start time:", cFmt(startTime))
	}

	totalElements := len(gtw.Config.Statics) + len(gtw.Config.Transitions)

	// Keep track of the total time. It is increased every time a new element duration is encountered.
	eventTime := startTime

	for i := 0; i < totalElements; i++ {
		// The duration of the event is specified in the XML file, but not when it should start

		// Get an element, by index. This is an interface{} and is expected to be a GStatic or a GTransition
		eInterface, err := gtw.Config.Get(i)
		if err != nil {
			return err
		}
		if s, ok := eInterface.(GStatic); ok {
			if verbose {
				fmt.Printf("Registering static event at %s for setting %s\n", cFmt(eventTime), s.Filename)
			}

			// Place values into variables, before enclosing it in the function below.
			from := eventTime
			window := s.Duration()
			cooldown := window
			imageFilename := s.Filename

			// Register a static event
			eventloop.Add(event.New(from, window, cooldown, func() {
				if verbose {
					fmt.Printf("Triggered static wallpaper event at %s\n", cFmt(from))
					fmt.Println("Window:", dFmt(window))
					fmt.Println("Cooldown:", dFmt(cooldown))
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
				fmt.Printf("Registering transition at %s for transitioning from %s to %s.\n", cFmt(eventTime), t.FromFilename, t.ToFilename)
			}

			// cross fade steps
			steps := 10

			from := eventTime
			window := t.Duration()
			upTo := eventTime.Add(window)
			cooldown := window / time.Duration(steps)
			tType := t.Type
			tFromFilename := t.FromFilename
			tToFilename := t.ToFilename
			loopWait := gtw.LoopWait
			//tempDir := ""
			var err error

			// Register a transition event
			eventloop.Add(event.New(from, window, cooldown, func() {
				progress := window - event.ToToday(upTo).Sub(event.ToToday(time.Now()))
				if progress < 0 {
					progress *= -1
				}
				ratio := float64(progress) / float64(window)
				if verbose {
					fmt.Printf("Triggered transition event at %s (%d%% complete)\n", cFmt(from), int(ratio*100))
					fmt.Println("Progress:", dFmt(progress))
					fmt.Println("Up to:", cFmt(upTo))
					fmt.Println("Window:", dFmt(window))
					fmt.Println("Cooldown:", dFmt(cooldown))
					fmt.Println("Loop wait:", dFmt(loopWait))
					fmt.Println("Transition type:", tType)
					fmt.Println("From filename", tFromFilename)
					fmt.Println("To filename", tToFilename)
				}

				//tempDir, err = ioutil.TempDir("", "crossfade")
				//if err != nil {
				//	fmt.Fprintf(os.Stderr, "Could not create temporary directory: %v\n", err)
				//	return // return from anon func
				//}
				tempDir := "/tmp"
				// Prepare to write an image to the temporary directory
				tempImageFilename := filepath.Join(tempDir, "gnometimed_crossfade.png") // .jpg is also possible

				//// Remove the temporary directory 5 minutes after this has passed
				//eventloop.Once(upTo.Add(5*time.Minute), func() {
				//	if exists(tempDir) {
				//		if verbose {
				//			fmt.Println("Removing", tempDir)
				//		}
				//		// Clean up
				//		os.RemoveAll(tempDir)
				//	}
				//})

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
	eventloop.Go(gtw.LoopWait)

	return nil
}
