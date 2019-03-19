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

	verbose := false
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		verbose = true
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Please specify an image filename.")
			os.Exit(1)
		}
		imageFilename = os.Args[2]
	}

	// Find the absolute path
	absImageFilename, err := filepath.Abs(imageFilename)
	if err == nil {
		imageFilename = absImageFilename
	}

	// Set the desktop wallpaper
	if err := wallutils.SetWallpaperVerbose(imageFilename, verbose); err != nil {
		fmt.Fprintf(os.Stderr, "Could not set wallpaper: %s\n", err)
		os.Exit(1)
	}
}
