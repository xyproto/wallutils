package simpletimed

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/xyproto/wallutils/pkg/event"
)

var setmut = &sync.RWMutex{}

// UntilNext finds the duration until the next event starts
func (stw *Wallpaper) UntilNext(et time.Time) (time.Duration, time.Time) {
	// Gather all start times from the list of transitions and list of static wallpaper commands
	var startTimes []time.Time
	for _, t := range stw.Transitions {
		startTimes = append(startTimes, t.From)
	}
	for _, s := range stw.Statics {
		startTimes = append(startTimes, s.At)
	}

	// Using all the collected hours&minutes, create a list of times both today and tomorrow that uses those hours&minutes
	now := time.Now()
	var allTimes []time.Time
	for _, t := range startTimes {
		today := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
		tomorrow := today.AddDate(0, 0, 1)
		allTimes = append(allTimes, today, tomorrow)
	}

	// Now we have all possible start times, now to find the ones that are both positive and smallest
	mindiff := h24
	when := now
	for _, t := range allTimes {
		diff := t.Sub(et)
		if diff > 0 && diff < mindiff {
			mindiff = diff
			when = t
		}
	}

	// Return the smallest time difference and the point in time for when that is
	return mindiff, when
}

// NextEvent finds the next event, given a timestamp.
// Returns an interface{} that is either a static or transition event.
func (stw *Wallpaper) NextEvent(et time.Time) (interface{}, time.Time, error) {
	// Create a map, from timestamps to wallpaper events
	events := make(map[time.Time]interface{})
	for _, t := range stw.Transitions {
		events[t.From] = t
	}
	for _, s := range stw.Statics {
		events[s.At] = s
	}
	if len(events) == 0 {
		return nil, et, errors.New("can not find next event: got no events")
	}

	// Using all the collected hours&minutes, create a list of times both today and tomorrow that uses those hours&minutes
	allTimes := make(map[time.Time]interface{})
	now := time.Now()
	for t, e := range events {
		today := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
		tomorrow := today.AddDate(0, 0, 1)
		allTimes[today] = e
		allTimes[tomorrow] = e
	}

	// Now we have all possible start times, now to find the ones that are both positive and smallest
	mindiff := h24
	when := now
	var eventHappening interface{}
	for t, e := range allTimes {
		if eventHappening == nil {
			eventHappening = e
		}
		diff := t.Sub(et)
		if diff > 0 && diff < mindiff {
			mindiff = diff
			when = t
			eventHappening = e
		}
	}

	return eventHappening, when, nil
}

// PrevEvent finds the previous event, given a timestamp.
// Returns an interface{} that is either a static or transition event.
func (stw *Wallpaper) PrevEvent(et time.Time) (interface{}, time.Time, error) {
	// Create a map, from timestamps to wallpaper events
	events := make(map[time.Time]interface{})
	for _, t := range stw.Transitions {
		events[t.From] = t
	}
	for _, s := range stw.Statics {
		events[s.At] = s
	}
	if len(events) == 0 {
		return nil, et, errors.New("can not find next event: got no events")
	}
	// Using all the collected hours&minutes, create a list of times both today and tomorrow that uses those hours&minutes
	allTimes := make(map[time.Time]interface{})
	now := time.Now()
	for t, e := range events {
		today := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
		tomorrow := today.AddDate(0, 0, 1)
		allTimes[today] = e
		allTimes[tomorrow] = e
	}

	// Now we have all possible start times, find the ones that are are below the given et,
	// but as small as possible.

	mindiff := h24
	when := now
	var eventHappening interface{}
	for t, e := range allTimes {
		if eventHappening == nil {
			eventHappening = e
		}
		diff := et.Sub(t) // reverse subtraction, to find the time comparison back in time
		if diff > 0 && diff < mindiff {
			mindiff = diff
			when = t
			eventHappening = e
		}
	}

	return eventHappening, when, nil
}

// SetInitialWallpaper will set the first wallpaper, before starting the event loop
func (stw *Wallpaper) SetInitialWallpaper(verbose bool, setWallpaperFunc func(string) error, tempImageFilename string) error {
	now := time.Now()
	e, whenPrev, err := stw.PrevEvent(now)
	if err != nil {
		return err
	}
	_, whenNext, err := stw.NextEvent(now)
	if err != nil {
		return err
	}

	// the length of the currently ongoing event
	eventLength := whenNext.Sub(whenPrev)

	switch v := e.(type) {
	case *Static:
		s := v

		// Place values into variables, before enclosing it in the function below.
		// from := s.At
		// elapsed := time.Now().Sub(when) // now - when the previous event was set to trigger
		// durationUntilNext, nextTime := stw.UntilNext(s.At)
		// window := mod24(durationUntilNext - elapsed) // duration until next event start, minus time elapsed
		// Duration until next event start, from now
		// window := time.Now().Sub(when)

		window := eventLength
		cooldown := eventLength

		imageFilename := s.Filename

		if verbose {
			fmt.Printf("Attaching to ongoing static wallpaper event that started at %s\n", cFmt(whenPrev))
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
			return fmt.Errorf("file does not exist: %s", imageFilename)
		}

		// Set the desktop wallpaper, if possible
		if verbose {
			fmt.Printf("Setting %s.\n", imageFilename)
		}
		if err := setWallpaperFunc(imageFilename); err != nil {
			return fmt.Errorf("could not set wallpaper: %v", err)
		}

		// Just sleep for half the cooldown, to have some time to register events too
		//if verbose {
		//fmt.Println("Activating events in", dFmt(cooldown/2))
		//}
		//time.Sleep(cooldown / 2)
	case *Transition:
		t := v

		elapsed := now.Sub(whenPrev)
		ratio := float64(elapsed) / float64(eventLength)

		from := whenPrev
		upTo := whenNext

		tType := t.Type
		tFromFilename := t.FromFilename
		loopWait := stw.LoopWait

		if verbose {
			fmt.Printf("Initial transition event at %s (%d%% complete)\n", cFmt(from), int(ratio*100))
			fmt.Println("Progress:", dFmt(elapsed))
			fmt.Println("Up to:", cFmt(upTo))
			fmt.Println("Window:", dFmt(eventLength))
			fmt.Println("Loop wait:", dFmt(loopWait))
			fmt.Println("Transition type:", tType)
			fmt.Println("Using filename", tFromFilename)
		}

		// Set the "from" image before crossfading, so that something happens immediately

		// Set the desktop wallpaper, if possible
		if verbose {
			fmt.Printf("Setting %s.\n", tFromFilename)
		}
		if err := setWallpaperFunc(tFromFilename); err != nil {
			return fmt.Errorf("could not set wallpaper: %v", err)
		}

	default:
		return errors.New("could not set initial wallpaper: no previous event")
	}
	return nil
}

// EventLoop will start the event loop for this Simple Timed Wallpaper
func (stw *Wallpaper) EventLoop(verbose bool, setWallpaperFunc func(string) error, tempImageFilename string) error {
	if verbose {
		fmt.Println("Using the Simple Timed Wallpaper format.")
	}

	// Listen for SIGHUP or SIGUSR1, to refresh the wallpaper.
	// Can be used after resume from sleep.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGUSR1)
	go func() {
		for {
			// Wait for a signal of the type given to signal.Notify
			sig := <-signals
			// Refresh the wallpaper
			fmt.Println("Received", sig)
			// Launch a goroutine for setting the wallpaper
			go func() {
				setmut.Lock()
				if err := stw.SetInitialWallpaper(verbose, setWallpaperFunc, tempImageFilename); err != nil {
					fmt.Fprintln(os.Stderr, "Error:", err)
				}
				setmut.Unlock()
			}()
		}
	}()

	setmut.Lock()
	if err := stw.SetInitialWallpaper(verbose, setWallpaperFunc, tempImageFilename); err != nil {
		setmut.Unlock()
		return err
	}
	setmut.Unlock()

	eventloop := event.NewSystem(stw.LoopWait)

	for _, s := range stw.Statics {
		if verbose {
			fmt.Printf("Event at %s for setting %s\n", cFmt(s.At), s.Filename)
		}

		// Place values into variables, before enclosing it in the function below.
		from := s.At

		nextEventDuration, _ := stw.UntilNext(s.At)

		window := mod24(nextEventDuration) // duration until next event start
		cooldown := window
		imageFilename := s.Filename

		// Register a static event
		eventloop.ClockEvent(from.Hour(), from.Minute(), func() error {
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
				return err // return from anon func
			}

			// Set the desktop wallpaper, if possible
			if verbose {
				fmt.Printf("Setting %s.\n", imageFilename)
			}
			if err := setWallpaperFunc(imageFilename); err != nil {
				fmt.Fprintf(os.Stderr, "Could not set wallpaper: %v\n", err)
				return err // return from anon func
			}
			return nil
		})
	}
	for _, t := range stw.Transitions {
		if verbose {
			fmt.Printf("Transition at %s from %s to %s.\n", cFmt(t.From), t.FromFilename, t.ToFilename)
		}

		from := t.From
		window := t.Duration()
		upTo := from.Add(window)
		tType := t.Type
		tFromFilename := t.FromFilename
		tToFilename := t.ToFilename
		loopWait := stw.LoopWait
		halfway := from.Add(window / 2)

		// Register the start of a transition event
		eventloop.ClockEvent(from.Hour(), from.Minute(), func() error {
			progress := mod24(window - event.ToToday(upTo).Sub(event.ToToday(time.Now())))
			ratio := float64(progress) / float64(window)
			if verbose {
				fmt.Printf("Triggered transition event at %s (%d%% complete)\n", cFmt(from), int(ratio*100))
				fmt.Println("Progress:", dFmt(progress))
				fmt.Println("Up to:", cFmt(upTo))
				fmt.Println("Window:", dFmt(window))
				fmt.Println("Loop wait:", dFmt(loopWait))
				fmt.Println("Transition type:", tType)
				fmt.Println("Using filename", tFromFilename)
			}
			tempImageFilename := tFromFilename
			// Set the desktop wallpaper, if possible
			if verbose {
				fmt.Printf("Setting %s.\n", tempImageFilename)
			}
			setmut.Lock()
			if err := setWallpaperFunc(tempImageFilename); err != nil {
				setmut.Unlock()
				fmt.Fprintf(os.Stderr, "Could not set wallpaper: %v\n", err)
				return err // return from anon func
			}
			setmut.Unlock()
			return nil
		})

		// Register a halfway transition event
		eventloop.ClockEvent(halfway.Hour(), halfway.Minute(), func() error {
			progress := mod24(window - event.ToToday(upTo).Sub(event.ToToday(time.Now())))
			ratio := float64(progress) / float64(window)
			if verbose {
				fmt.Printf("Triggered transition event at %s (%d%% complete)\n", cFmt(from), int(ratio*100))
				fmt.Println("Progress:", dFmt(progress))
				fmt.Println("Up to:", cFmt(upTo))
				fmt.Println("Window:", dFmt(window))
				// fmt.Println("Cooldown:", dFmt(cooldown))
				fmt.Println("Loop wait:", dFmt(loopWait))
				fmt.Println("Transition type:", tType)
				fmt.Println("From filename", tFromFilename)
				fmt.Println("To filename", tToFilename)
				fmt.Println("Crossfading between images.")
			}
			// Crossfade and write the new image to the temporary directory
			tFromImg, err := imgio.Open(tFromFilename)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return err
			}
			tToImg, err := imgio.Open(tToFilename)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return err
			}
			// Crossfade and write the new image to the temporary directory
			setmut.Lock()
			blendedImage := blend.Opacity(tFromImg, tToImg, ratio)
			err = imgio.Save(tempImageFilename, blendedImage, imgio.JPEGEncoder(100))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not crossfade images in transition: %v\n", err)
				setmut.Unlock()
				return err
			}
			setmut.Unlock()
			// Double check that the generated file exists
			if _, err := os.Stat(tempImageFilename); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "File does not exist: %s\n", tempImageFilename)
				return err // return from anon func
			}
			// Set the desktop wallpaper, if possible
			if verbose {
				fmt.Printf("Setting %s.\n", tempImageFilename)
			}
			setmut.Lock()
			if err := setWallpaperFunc(tempImageFilename); err != nil {
				setmut.Unlock()
				fmt.Fprintf(os.Stderr, "Could not set wallpaper: %v\n", err)
				return err // return from anon func
			}
			setmut.Unlock()
			return nil
		})

	}

	// Endless loop! Will wait LoopWait duration between each event loop cycle.
	eventloop.Run(verbose)

	// eventloop.Run returns immediately, so start an endless loop
	for {
		time.Sleep(1 * time.Second)
	}

	return nil
}
