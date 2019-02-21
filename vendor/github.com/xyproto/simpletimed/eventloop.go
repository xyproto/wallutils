package simpletimed

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/xyproto/crossfade"
	"github.com/xyproto/event"
)

// UntilNext finds the duration until the next event starts
func (stw *Wallpaper) UntilNext(et time.Time) time.Duration {
	var startTimes []time.Time
	for _, t := range stw.Transitions {
		startTimes = append(startTimes, t.From)
	}
	for _, s := range stw.Statics {
		startTimes = append(startTimes, s.At)
	}
	mindiff := h24
	// OK, have all start times, now to find the ones that are both positive and smallest
	for _, st := range startTimes {
		//diff := st.Sub(et)
		diff := mod24(event.ToToday(et).Sub(event.ToToday(st)))
		if diff > 0 && diff < mindiff {
			mindiff = diff
		}
	}
	return mindiff
}

// NextEvent finds the next event, given a timestamp.
// Returns an interface{} that is either a static or transition event.
func (stw *Wallpaper) NextEvent(now time.Time) (interface{}, error) {
	// Create a map, from timestamps to wallpaper events
	events := make(map[time.Time]interface{})
	for _, t := range stw.Transitions {
		events[t.From] = t
	}
	for _, s := range stw.Statics {
		events[s.At] = s
	}
	if len(events) == 0 {
		return nil, errors.New("can not find next event: got no events")
	}
	// Go though all the event time stamps, and find the one that has the smallest (now time - event time)
	minDiff := h24
	var minEvent interface{}
	for t, e := range events {
		//fmt.Printf("now is: %v (%T)\n", now, now)
		//fmt.Printf("t is: %v (%T)\n", t, t)
		diff := mod24(event.ToToday(t).Sub(event.ToToday(now)))
		//fmt.Println("Diff for", c(t), ":", diff)
		if diff > 0 && diff < minDiff {
			minDiff = diff
			minEvent = e
			//fmt.Println("NEW SMALLEST DIFF RIGHT AFTER", c(now), c(t), minDiff)
		}
	}
	return minEvent, nil
}

// PrevEvent finds the previous event, given a timestamp.
// Returns an interface{} that is either a static or transition event.
func (stw *Wallpaper) PrevEvent(now time.Time) (interface{}, error) {
	// Create a map, from timestamps to wallpaper events
	events := make(map[time.Time]interface{})
	for _, t := range stw.Transitions {
		events[t.From] = t
	}
	for _, s := range stw.Statics {
		events[s.At] = s
	}
	if len(events) == 0 {
		return nil, errors.New("can not find previous event: got no events")
	}
	// Go though all the event time stamps, and find the one that has the smallest (now time - event time)
	minDiff := h24
	var minEvent interface{}
	for t, e := range events {
		//fmt.Printf("now is: %v (%T)\n", now, now)
		//fmt.Printf("t is: %v (%T)\n", t, t)
		diff := mod24(event.ToToday(now).Sub(event.ToToday(t)))
		//fmt.Println("Diff for", c(t), ":", diff)
		if diff > 0 && diff < minDiff {
			minDiff = diff
			minEvent = e
			//fmt.Println("NEW SMALLEST DIFF RIGHT BEFORE", c(now), c(t), minDiff)
		}
	}
	return minEvent, nil
}

// SetInitialWallpaper will set the first wallpaper, before starting the event loop
func (stw *Wallpaper) SetInitialWallpaper(verbose bool, setWallpaperFunc func(string) error) error {
	e, err := stw.PrevEvent(time.Now())
	if err != nil {
		return err
	}
	switch v := e.(type) {
	case *Static:
		s := v

		// Place values into variables, before enclosing it in the function below.
		from := s.At
		//elapsed := time.Now().Sub(s.At)
		elapsed := event.ToToday(time.Now()).Sub(event.ToToday(s.At))
		window := mod24(stw.UntilNext(s.At) - elapsed) // duration until next event start, minus time elapsed
		cooldown := window

		imageFilename := s.Filename

		if verbose {
			fmt.Printf("Initial static wallpaper event at %s\n", c(from))
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
			return fmt.Errorf("File does not exist: %s\n", imageFilename)
		}

		// Set the desktop wallpaper, if possible
		if verbose {
			fmt.Printf("Setting %s.\n", imageFilename)
		}
		if err := setWallpaperFunc(imageFilename); err != nil {
			return fmt.Errorf("Could not set wallpaper: %v\n", err)
		}

		// Just sleep for half the cooldown, to have some time to register events too
		if verbose {
			fmt.Println("Activating events in", cooldown/2)
		}
		time.Sleep(cooldown / 2)
	case *Transition:
		t := v

		now := time.Now()
		window := t.Duration()
		progress := mod24(window - event.ToToday(t.UpTo).Sub(event.ToToday(now)))
		ratio := float64(progress) / float64(window)
		from := t.From
		steps := 10
		cooldown := window / time.Duration(steps)
		upTo := from.Add(window)
		tType := t.Type
		tFromFilename := t.FromFilename
		tToFilename := t.ToFilename
		loopWait := stw.LoopWait
		tempDir := ""
		var err error

		if verbose {
			fmt.Printf("Initial transition event at %s (%d%% complete)\n", c(from), int(ratio*100))
			fmt.Println("Progress:", progress)
			fmt.Println("Up to:", c(upTo))
			fmt.Println("Window:", window)
			fmt.Println("Cooldown:", cooldown)
			fmt.Println("Loop wait:", loopWait)
			fmt.Println("Transition type:", tType)
			fmt.Println("From filename", tFromFilename)
			fmt.Println("To filename", tToFilename)
		}

		tempDir, err = ioutil.TempDir("", "crossfade")
		if err != nil {
			return fmt.Errorf("Could not create temporary directory: %v\n", err)
		}
		// Prepare to write an image to the temporary directory
		tempImageFilename := filepath.Join(tempDir, "out.png") // .jpg is also possible

		// Set the "from" image before crossfading, so that something happens immediately

		// Set the desktop wallpaper, if possible
		if verbose {
			fmt.Printf("Setting %s.\n", tFromFilename)
		}
		if err := setWallpaperFunc(tFromFilename); err != nil {
			return fmt.Errorf("Could not set wallpaper: %v\n", err)
		}

		if verbose {
			fmt.Println("Crossfading between images.")
		}

		// Crossfade and write the new image to the temporary directory
		if crossfade.Files(tFromFilename, tToFilename, tempImageFilename, ratio) != nil {
			return fmt.Errorf("Could not crossfade images in transition: %v\n", err)
		}

		// Double check that the generated file exists
		if _, err := os.Stat(tempImageFilename); os.IsNotExist(err) {
			return fmt.Errorf("File does not exist: %s\n", tempImageFilename)
		}

		// Set the desktop wallpaper, if possible
		if verbose {
			fmt.Printf("Setting %s.\n", tempImageFilename)
		}
		if err := setWallpaperFunc(tempImageFilename); err != nil {
			return fmt.Errorf("Could not set wallpaper: %v\n", err)
		}

		// Just sleep for half the cooldown, to have some time to register events too
		if verbose {
			fmt.Println("Activating events in", cooldown/2)
		}
		time.Sleep(cooldown / 2)

		// Remove the temporary directory 5 minutes after this
		go func() {
			time.Sleep(5 * time.Minute)
			if exists(tempDir) {
				if verbose {
					fmt.Println("Removing", tempDir)
				}
				// Clean up
				os.RemoveAll(tempDir)
			}
		}()

	default:
		return errors.New("could not set initial wallpaper: no previous event")
	}
	return nil
}

// EventLoop will start the event loop for this Simple Timed Wallpaper
func (stw *Wallpaper) EventLoop(verbose bool, setWallpaperFunc func(string) error) error {

	if verbose {
		fmt.Println("Using the Simple Timed Wallpaper format.")
	}

	stw.SetInitialWallpaper(verbose, setWallpaperFunc)

	eventloop := event.NewLoop()

	for _, s := range stw.Statics {
		if verbose {
			fmt.Printf("Registering static event at %s for setting %s\n", c(s.At), s.Filename)
		}

		// Place values into variables, before enclosing it in the function below.
		from := s.At
		window := mod24(stw.UntilNext(s.At)) // duration until next event start
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
	}
	for _, t := range stw.Transitions {
		if verbose {
			fmt.Printf("Registering transition at %s for transitioning from %s to %s.\n", c(t.From), t.FromFilename, t.ToFilename)
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
		loopWait := stw.LoopWait
		tempDir := ""
		var err error

		// Register a transition event
		//eventloop.Add(event.New(from, window, cooldown, event.ProgressWrapperInterval(from, upTo, loopWait, func(ratio float64) {
		eventloop.Add(event.New(from, window, cooldown, func() {
			progress := mod24(window - event.ToToday(upTo).Sub(event.ToToday(time.Now())))
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

			tempDir, err = ioutil.TempDir("", "crossfade")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not create temporary directory: %v\n", err)
				return // return from anon func
			}

			// Prepare to write an image to the temporary directory
			tempImageFilename := filepath.Join(tempDir, "out.png") // .jpg is also possible

			// Remove the temporary directory 5 minutes after this event has passed
			eventloop.Once(upTo.Add(5*time.Minute), func() {
				if exists(tempDir) {
					if verbose {
						fmt.Println("Removing", tempDir)
					}
					// Clean up
					os.RemoveAll(tempDir)
				}
			})

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
	}

	// Endless loop! Will wait LoopWait duration between each event loop cycle.
	eventloop.Go(stw.LoopWait)
	return nil
}
