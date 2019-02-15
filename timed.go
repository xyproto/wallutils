package monitor

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

// SetInitialWallpaper will set the first wallpaper, before starting the vent loop
func (stw *SimpleTimedWallpaper) SetInitialWallpaper(verbose bool) error {
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
		window := stw.UntilNext(s.At) - elapsed // duration until next event start, minus time elapsed
		for window < 0 {
			window += h24
		}
		for window > h24 {
			window -= h24
		}
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
		if err := SetWallpaperVerbose(imageFilename, verbose); err != nil {
			return fmt.Errorf("Could not set wallpaper: %v\n", err)
		}

		// Just sleep for half the cooldown, to have some time to register events too
		fmt.Println("Sleeping for", cooldown/2)
		time.Sleep(cooldown / 2)
	case *Transition:
		t := v

		now := time.Now()
		window := t.Duration()
		progress := window - event.ToToday(t.UpTo).Sub(event.ToToday(now))
		for progress > h24 {
			progress -= h24
		}
		for progress < 0 {
			progress += h24
		}
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

		if exists(tempDir) {
			// Clean up
			os.RemoveAll(tempDir)
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
		if err := SetWallpaperVerbose(tFromFilename, verbose); err != nil {
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
		if err := SetWallpaperVerbose(tempImageFilename, verbose); err != nil {
			return fmt.Errorf("Could not set wallpaper: %v\n", err)
		}

		// Just sleep for half the cooldown, to have some time to register events too
		fmt.Println("Sleeping for", cooldown/2)
		time.Sleep(cooldown / 2)
	default:
		return errors.New("could not set initial wallpaper: no previous event")
	}
	return nil
}

// EventLoop will start the event loop for this Simple Timed Wallpaper
func (stw *SimpleTimedWallpaper) EventLoop(verbose bool) error {

	if verbose {
		fmt.Println("Using the Simple Timed Wallpaper format.")
	}

	stw.SetInitialWallpaper(verbose)

	eventloop := event.NewLoop()

	for _, s := range stw.Statics {
		if verbose {
			fmt.Printf("Registering static event at %s for setting %s\n", c(s.At), s.Filename)
		}

		// Place values into variables, before enclosing it in the function below.
		from := s.At
		window := stw.UntilNext(s.At) // duration until next event start
		for window < 0 {
			window += h24
		}
		for window > h24 {
			window -= h24
		}
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
			if err := SetWallpaperVerbose(imageFilename, verbose); err != nil {
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
			if err := SetWallpaperVerbose(tempImageFilename, verbose); err != nil {
				fmt.Fprintf(os.Stderr, "Could not set wallpaper: %v\n", err)
				return // return from anon func
			}
		}))
	}

	// Endless loop! Will wait LoopWait duration between each event loop cycle.
	eventloop.Go(stw.LoopWait)
	return nil
}

// EventLoop will start the event loop for this GNOME Timed Wallpaper
func (gw *GnomeTimedWallpaper) EventLoop(verbose bool) error {

	if verbose {
		fmt.Println("Using the GNOME Timed Wallpaper format")
	}

	// Convert to a SimpleTimedWallpaper, only for setting the initial wallpaper
	stw, err := GnomeToSimple(gw)
	if err != nil {
		return err
	}
	stw.SetInitialWallpaper(verbose)

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
				if err := SetWallpaperVerbose(imageFilename, verbose); err != nil {
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
				if err := SetWallpaperVerbose(tempImageFilename, verbose); err != nil {
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

// UntilNext finds the duration until the next event starts
func (stw *SimpleTimedWallpaper) UntilNext(et time.Time) time.Duration {
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
		diff := event.ToToday(et).Sub(event.ToToday(st))
		for diff < 0 {
			diff += h24
		}
		for diff > h24 {
			diff -= h24
		}
		if diff > 0 && diff < mindiff {
			mindiff = diff
		}
	}
	return mindiff
}

// NextEvent finds the next event, given a timestamp.
// Returns an interface{} that is either a static or transition event.
func (stw *SimpleTimedWallpaper) NextEvent(now time.Time) (interface{}, error) {
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
		diff := event.ToToday(t).Sub(event.ToToday(now))
		//diff := t.Sub(now)
		for diff < 0 {
			diff += h24
		}
		for diff > h24 {
			diff -= h24
		}
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
func (stw *SimpleTimedWallpaper) PrevEvent(now time.Time) (interface{}, error) {
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
		diff := event.ToToday(now).Sub(event.ToToday(t))
		//diff := now.Sub(t)
		for diff < 0 {
			diff += h24
		}
		for diff > h24 {
			diff -= h24
		}
		//fmt.Println("Diff for", c(t), ":", diff)
		if diff > 0 && diff < minDiff {
			minDiff = diff
			minEvent = e
			//fmt.Println("NEW SMALLEST DIFF RIGHT BEFORE", c(now), c(t), minDiff)
		}
	}
	return minEvent, nil
}
