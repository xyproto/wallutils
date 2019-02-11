package main

import (
	"fmt"
	"os"

	"github.com/xyproto/monitor"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "Please give a timed wallpaper name as the first argument.")
		os.Exit(1)
	}

	collectionName := os.Args[1]
	fmt.Printf("Setting wallpaper collection \"%s\"\n", collectionName)

	fmt.Print("Searching for wallpapers...")
	_, gnomeWallpapers, simpleTimedWallpapers := monitor.FindWallpapers()
	if len(gnomeWallpapers) == 0 && len(simpleTimedWallpapers) == 0 {
		fmt.Fprintln(os.Stderr, "Could not find any timed wallpapers on the system.")
		os.Exit(1)
	} else {
		fmt.Println("ok")
	}

	fmt.Print("Filtering wallpapers by name...")
	simpleTimedWallpapers = monitor.FilterSimpleTimedWallpapers(collectionName, simpleTimedWallpapers)
	gnomeWallpapers = monitor.FilterGnomeWallpapers(collectionName, gnomeWallpapers)
	fmt.Println("ok")

	// gnomeWallpapers and simpleTimedWallpapers have now been filtered so that they only contain elements with matching collection names

	if (len(gnomeWallpapers) == 0) && (len(simpleTimedWallpapers) == 0) {
		fmt.Fprintln(os.Stderr, "No such timed wallpaper: "+collectionName)
		os.Exit(1)
	}

	if len(simpleTimedWallpapers) == 1 {
		err := monitor.SetSimpleTimedWallpaper(simpleTimedWallpapers[0], true)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else if len(gnomeWallpapers) == 1 {
		err := monitor.SetGnomeTimedWallpaper(gnomeWallpapers[0], true)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else if (len(gnomeWallpapers) > 1) || (len(simpleTimedWallpapers) > 1) {
		fmt.Fprintln(os.Stderr, "Found several timed backgrounds, with the same name.")
		os.Exit(1)
	}
}
