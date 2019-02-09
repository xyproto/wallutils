package main

import (
	"fmt"
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

func setGnomeWallpaper(gw *monitor.GnomeWallpaper) error {

	fmt.Println("--- setGnomeWallpaper ---")
	fmt.Println("collection name:", gw.CollectionName)

	eventloop := event.NewLoop()

	// Get the start time for the wallpaper collection (which is offset by X
	// seconds per static wallpaper)
	startTime := gw.StartTime()

	fmt.Println("Timed wallpaper start time:", startTime)

	totalElements := len(gw.Config.Statics) + len(gw.Config.Transitions)

	// Keep track of the total time. It is increased every time a new element duration is encountered.
	eventTime := startTime

	// The cooldown for every event is 15 minutes. It can not be retriggered in that time period.
	cooldown := 15 * time.Minute

	for i := 0; i < totalElements; i++ {
		// The duration of the event is specified in the XML file, but not when it should start
		var window time.Duration

		eInterface, err := gw.Config.Get(i)
		if err != nil {
			fmt.Println("GAH!", err)
			break
		}
		if s, ok := eInterface.(monitor.GStatic); ok {
			//fmt.Println("GOT STATIC", s)
			window = s.Duration()
			fmt.Println("WINDOW", window)
			fmt.Println("EVENT TIME", eventTime)

			// Place the filename into a variable, before enclosing it in the
			// function below.
			imageFilename := s.Filename

			fmt.Printf("Registering event at %s for changing to %s\n", eventTime, imageFilename)
			eventloop.Add(event.New(eventTime, window, cooldown, func() {
				fmt.Println("WALLPAPER EVENT", imageFilename)

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

			fmt.Println("EVENT START", eventTime)
			// Increase the variable that keeps track of the time!
			eventTime = eventTime.Add(window)
			fmt.Println("EVENT END", eventTime)
		} else if t, ok := eInterface.(monitor.GTransition); ok {
			//fmt.Println("GOT TRANSITION", t)
			window = t.Duration()
			fmt.Println("WINDOW", window)

			fmt.Println("TYPE         ", t.Type)
			fmt.Println("FROM FILENAME", t.FromFilename)
			fmt.Println("TO FILENAME  ", t.ToFilename)

			fmt.Println("!!!TO IMPLEMENT!!!")

			fmt.Println("EVENT START", eventTime)
			// Increase the variable that keeps track of the time!
			eventTime = eventTime.Add(window)
			fmt.Println("EVENT END", eventTime)
		} else {
			fmt.Println("GOT NOTHING")
			break
		}

	}

	// QUIT
	return nil

	// Endless loop! Will wait 5 seconds between each event loop cycle
	eventloop.Go(5 * time.Second)

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
