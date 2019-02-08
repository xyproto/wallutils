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

func setGnomeWallpaper(gw *monitor.GnomeWallpaper) error {

	events := event.NewEvents()

	fmt.Println("--- setGnomeWallpaper ---")
	fmt.Println("collection name:", gw.CollectionName)

	// Plan:
	// List all "static" wallpaper filenames, with associated starttime + duration

	// Get the base start time for today
	timedWallpaperStart := gw.StartTimeToday()

	fmt.Println("start time:", gw.Time())
	for _, s := range gw.Config.Statics {
		seconds := time.Duration(s.Seconds) * time.Second
		eventStart := timedWallpaperStart.Add(seconds)

		// TODO: Add support for hour/minute/second events in the even package.
		//       For now, these events only work for one day.

		fmt.Println("trigger time", triggerStart.Add(seconds))
		fn := s.Filename
		fmt.Println("image", fn, seconds)

		events.Add

	}

	fmt.Println("start time:", gw.Time())

	fmt.Println("TO IMPLEMENT: GNOME TIMED BACKGROUND")
	os.Exit(1)
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
