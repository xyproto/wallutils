package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xyproto/wallutils"
)

const versionString = "setwallpaper"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println(wallutils.VersionString)
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Please specify an image filename.")
		os.Exit(1)
	}
	imageFilename := os.Args[1]

	// Find the absolute path
	absImageFilename, err := filepath.Abs(imageFilename)
	if err == nil {
		imageFilename = absImageFilename
	}

	// Check that the file exists
	if _, err := os.Stat(imageFilename); os.IsNotExist(err) {
		// File does not exist
		fmt.Fprintf(os.Stderr, "File does not exist: %s\n", imageFilename)
		os.Exit(1)
	}

	// Set the desktop wallpaper
	if err := wallutils.SetWallpaper(imageFilename); err != nil {
		fmt.Fprintf(os.Stderr, "Could not set wallpaper: %s\n", err)
		os.Exit(1)
	}
}
