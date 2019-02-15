package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/xyproto/monitor"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println(monitor.VersionString)
		os.Exit(0)
	}

	if len(os.Args) == 1 {
		fmt.Println("Please specify a directory to choose wallpapers from.")
		os.Exit(1)
	}

	dir := os.Args[1]

	pngMatches, err := filepath.Glob(filepath.Join(dir, "/*.png"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	jpgMatches, err := filepath.Glob(filepath.Join(dir, "/*.jpg"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	matches := append(pngMatches, jpgMatches...)

	if len(matches) == 0 {
		fmt.Fprintln(os.Stderr, "Found no png or jpg files in "+dir)
		os.Exit(1)
	}

	imageFilename := matches[rand.Int()%len(matches)]
	if absImageFilename, err := filepath.Abs(imageFilename); err == nil {
		imageFilename = absImageFilename
	}
	if _, err := os.Stat(imageFilename); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Setting background image to: " + imageFilename)
	if err := monitor.SetWallpaper(imageFilename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
