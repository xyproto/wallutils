package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xyproto/wallutils"
)

// Select the wallpaper that is closest to the current monitor resolution and set that as the wallpaper
func SelectAndSetWallpaper(wallpapers []*wallutils.Wallpaper) error {
	// Gather a slice of filenames
	var filenames []string
	for _, wp := range wallpapers {
		filenames = append(filenames, wp.Path)
	}

	// Select the image that is closest to the monitor resolution
	imageFilename, err := wallutils.Closest(filenames)
	if err != nil {
		return err
	}

	// Find the absolute path
	absImageFilename, err := filepath.Abs(imageFilename)
	if err == nil {
		imageFilename = absImageFilename
	}

	// Set the desktop wallpaper
	if err := wallutils.SetWallpaper(imageFilename); err != nil {
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
	searchResults, err := wallutils.FindWallpapers()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	wallpapers := searchResults.Wallpapers()
	gnomeTimedWallpapers := searchResults.GnomeTimedWallpapers()
	simpleTimedWallpapers := searchResults.SimpleTimedWallpapers()

	if searchResults.Empty() {
		fmt.Fprintln(os.Stderr, "Could not find any wallpapers on the system.")
		os.Exit(1)
	} else {
		fmt.Println("ok")
	}

	fmt.Print("Filtering wallpapers by collection name...")
	wallpapers = searchResults.WallpapersByName(collectionName)
	gnomeTimedWallpapers = searchResults.GnomeTimedWallpapersByName(collectionName)
	simpleTimedWallpapers = searchResults.SimpleTimedWallpapersByName(collectionName)
	fmt.Println("ok")

	if len(wallpapers) == 0 && (len(gnomeTimedWallpapers) > 0 || len(simpleTimedWallpapers) > 0) {
		fmt.Fprintln(os.Stderr, "Timed wallpapers are not supported by this utility.")
		os.Exit(1)
	}

	if len(wallpapers) == 0 {
		fmt.Fprintln(os.Stderr, "No such collection: "+collectionName)
		os.Exit(1)
	}

	if err = SelectAndSetWallpaper(wallpapers); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
