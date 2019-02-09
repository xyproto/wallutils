package main

import (
	"fmt"
	"os"

	"github.com/xyproto/monitor"
)

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

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "Please give a timed wallpaper name as the first argument.")
		os.Exit(1)
	}

	collectionName := os.Args[1]
	fmt.Printf("Setting wallpaper collection \"%s\"\n", collectionName)

	fmt.Print("Searching for wallpapers...")
	_, gnomeWallpapers := monitor.FindWallpapers()
	if len(gnomeWallpapers) == 0 {
		fmt.Fprintln(os.Stderr, "Could not find any timed wallpapers on the system.")
		os.Exit(1)
	} else {
		fmt.Println("ok")
	}

	fmt.Print("Filtering wallpapers by name...")
	gnomeWallpapers = filterGnomeWallpapers(collectionName, gnomeWallpapers)
	fmt.Println("ok")

	if len(gnomeWallpapers) == 0 {
		fmt.Fprintln(os.Stderr, "No such tiemd wallpaper: "+collectionName)
		os.Exit(1)
	}

	// gnomeWallpapers are now filtered to only contain elements with matching collection names

	if len(gnomeWallpapers) == 1 {
		err := monitor.SetTimedWallpaper(gnomeWallpapers[0], true)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else if len(gnomeWallpapers) > 1 {
		fmt.Fprintln(os.Stderr, "Found several GNOME timed backgrounds, with the same name.")
		os.Exit(1)
	}
}
