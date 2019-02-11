package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xyproto/monitor"
)

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

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "Please give a wallpaper collection name as the first argument.")
		os.Exit(1)
	}

	collectionName := os.Args[1]
	fmt.Printf("Setting wallpaper collection \"%s\"\n", collectionName)

	fmt.Print("Searching for wallpapers...")
	wallpapers, gnomeWallpapers, simpleTimedWallpapers := monitor.FindWallpapers()
	if len(wallpapers) == 0 && len(gnomeWallpapers) == 0 && len(simpleTimedWallpapers) == 0 {
		fmt.Fprintln(os.Stderr, "Could not find any wallpapers on the system.")
		os.Exit(1)
	} else {
		fmt.Println("ok")
	}

	fmt.Print("Filtering wallpapers by collection name...")
	wallpapers = monitor.FilterWallpapers(collectionName, wallpapers)
	gnomeWallpapers = monitor.FilterGnomeWallpapers(collectionName, gnomeWallpapers)
	simpleTimedWallpapers = monitor.FilterSimpleTimedWallpapers(collectionName, simpleTimedWallpapers)
	fmt.Println("ok")

	if len(wallpapers) == 0 && len(gnomeWallpapers) == 0 && len(simpleTimedWallpapers) == 0 {
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
		fmt.Fprintln(os.Stderr, "Timed wallpapers are not supported by this utility")
		os.Exit(1)
	} else if len(wallpapers) == 0 && len(gnomeWallpapers) > 1 {
		fmt.Fprintln(os.Stderr, "Found several GNOME timed backgrounds, with the same collection name!")
		os.Exit(1)
	}
}
