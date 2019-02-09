package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/xyproto/event"
	"github.com/xyproto/monitor"
)

// filterWallpapers will filter out wallpapers that both match with the collection name, and are also marked as part of a collection
func filterWallpapers(collectionName string, wallpapers []*monitor.Wallpaper) []*monitor.Wallpaper {
	var collection []*monitor.Wallpaper
	for _, wp := range wallpapers {
		if wp.PartOfCollection && wp.CollectionName == collectionName {
			collection = append(collection, wp)
		}
	}
	return collection
}

// filterGnomeWallpapers will filter out gnome timed wallpapers that match with the collection name
func filterGnomeWallpapers(collectionName string, gnomeWallpapers []*monitor.GnomeWallpaper) []*monitor.GnomeWallpaper {
	var collection []*monitor.GnomeWallpaper
	for _, gw := range gnomeWallpapers {
		if gw.CollectionName == collectionName {
			collection = append(collection, gw)
		}
	}
	return collection
}

func setWallpaper(wallpapers []*monitor.Wallpaper) error {
	// Gather a slice of filenames
	var filenames []string
	for _, wp := range wallpapers {
		filenames = append(filenames, wp.Path)
	}

	// Select the image that is closest to the monitor resolution
	imageFilename, err := monitor.Closest(filenames)
	if err != nil {
		return err
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

	// Set the desktop wallpaper
	if err := monitor.SetWallpaper(imageFilename); err != nil {
		return fmt.Errorf("Could not set wallpaper: %s\n", err)
	}

	return nil
}

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

func c(t time.Time) string {
	return fmt.Sprintf("%.2d:%.2d", t.Hour(), t.Minute())
}

func setGnomeWallpaper(gw *monitor.GnomeWallpaper) error {

	loopWait := 5 * time.Second

	fmt.Println("Using timed wallpaper:", gw.CollectionName)

	eventloop := event.NewLoop()

	// Get the start time for the wallpaper collection (which is offset by X
	// seconds per static wallpaper)
	startTime := gw.StartTime()

	// The start time of the timed wallpaper as a whole
	fmt.Println("Timed wallpaper start time:", c(startTime))

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
		if s, ok := eInterface.(monitor.GStatic); ok {
			window := s.Duration()
			// We have a static GNOME wallpaper element, with a duration and an image filename

			fmt.Printf("Registering STATIC at %s for changing to %s\n", c(eventTime), s.Filename)

			// Place values into variables, before enclosing it in the function below.
			from := eventTime
			cooldown := window
			//imageFilename := s.Filename
			eventloop.Add(event.New(from, window, cooldown, func() {
				fmt.Println("TRIGGERED STATIC WALLPAPER EVENT")

				fmt.Println("FROM", c(from))
				fmt.Println("WINDOW", window)
				fmt.Println("COOLDOWN", cooldown)

				// Enclose variable
				imageFilename := s.Filename
				fmt.Println("FILENAME", imageFilename)

				// Find the absolute path
				absImageFilename, err := filepath.Abs(imageFilename)
				if err == nil {
					imageFilename = absImageFilename
				}

				// Check that the file exists
				if _, err := os.Stat(imageFilename); os.IsNotExist(err) {
					fmt.Errorf("File does not exist: %s\n", imageFilename)
					return // return from anon func
				}

				// Set the desktop wallpaper, if possible
				if err := monitor.SetWallpaper(imageFilename); err != nil {
					fmt.Errorf("Could not set wallpaper: %s\n", err)
					return // return from anon func
				}
			}))
			// Increase the variable that keeps track of the time
			//eventTime = event.ToToday(eventTime.Add(window))
			eventTime = eventTime.Add(window)

		} else if t, ok := eInterface.(monitor.GTransition); ok {
			// Increase the variable that keeps track of the time
			window := t.Duration()

			// We have a GNOME wallpaper transition, with a duration, a type,
			// and two image filenames.

			fmt.Printf("Registering TRANSITION at %s for transitioning to %s.\n", c(eventTime), t.ToFilename)

			from := eventTime
			cooldown := window
			upTo := eventTime.Add(window)
			eventloop.Add(event.New(from, window, cooldown, event.ProgressWrapperInterval(from, upTo, loopWait, func(p float64) {

				fmt.Println("TRIGGERED TRANSITION EVENT")

				fmt.Println("TO IMPLEMENT: A smooth transition")

				// Enclose variables
				tType := t.Type
				tFromFilename := t.FromFilename
				tToFilename := t.ToFilename

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

				// TODO: Create a temporary image that is a mix between t.FromFilename and t.ToFilename, using p as the ratio
				imageFilename := tToFilename

				// Find the absolute path
				absImageFilename, err := filepath.Abs(imageFilename)
				if err == nil {
					imageFilename = absImageFilename
				}

				// Check that the file exists
				if _, err := os.Stat(imageFilename); os.IsNotExist(err) {
					fmt.Errorf("File does not exist: %s\n", imageFilename)
					return // return from anon func
				}

				// Set the desktop wallpaper, if possible
				if err := monitor.SetWallpaper(imageFilename); err != nil {
					fmt.Errorf("Could not set wallpaper: %s\n", err)
					return // return from anon func
				}
			})))

			//eventTime = event.ToToday(eventTime.Add(window))
			eventTime = eventTime.Add(window)
		} else {
			log.Println("warning: got an element that is not a monitor.GStatic and not a monitor.GTransition")
			continue
		}
	}

	// Endless loop! Will wait loopWait duration between each event loop cycle
	eventloop.Go(loopWait)

	return nil
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "Please give a wallpaper collection name as the first argument.")
		os.Exit(1)
	}

	collectionName := os.Args[1]
	fmt.Printf("Setting wallpaper collection \"%s\"\n", collectionName)

	fmt.Print("Searching for wallpapers...")
	wallpapers, gnomeWallpapers := monitor.FindWallpapers()
	if len(wallpapers) == 0 && len(gnomeWallpapers) == 0 {
		fmt.Fprintln(os.Stderr, "Could not find any wallpapers on the system.")
		os.Exit(1)
	} else {
		fmt.Println("ok")
	}

	fmt.Print("Filtering wallpapers by collection name...")
	wallpapers = filterWallpapers(collectionName, wallpapers)
	gnomeWallpapers = filterGnomeWallpapers(collectionName, gnomeWallpapers)
	fmt.Println("ok")

	if len(wallpapers) == 0 && len(gnomeWallpapers) == 0 {
		fmt.Fprintln(os.Stderr, "No such collection: "+collectionName)
		os.Exit(1)
	}

	// wallpapers and gnomeWallpapers are now filtered to only contain elements with matching collection names

	if len(wallpapers) > 0 && len(gnomeWallpapers) > 0 {
		fmt.Fprintln(os.Stderr, "Found both a wallpaper collection and a GNOME timed background after filtering by collection name.")
		os.Exit(1)
	}
	if len(wallpapers) > 0 && len(gnomeWallpapers) == 0 {
		err := setWallpaper(wallpapers)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else if len(wallpapers) == 0 && len(gnomeWallpapers) == 1 {
		err := setGnomeWallpaper(gnomeWallpapers[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else if len(wallpapers) == 0 && len(gnomeWallpapers) > 1 {
		fmt.Fprintln(os.Stderr, "Found several GNOME timed backgrounds, with the same collection name!")
		os.Exit(1)
	}
}
