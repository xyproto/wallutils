package main

import (
	"fmt"
	"github.com/xyproto/monitor"
	"os"
	"path/filepath"
	"text/tabwriter"
)

func has(sl []string, s string) bool {
	for _, e := range sl {
		if e == s {
			return true
		}
	}
	return false
}

func main() {
	alsoPrintPath := len(os.Args) > 1 && os.Args[1] == "-l"

	if !alsoPrintPath {
		for _, name := range monitor.FindCollectionNames() {
			fmt.Println(name)
		}
		return
	}

	// Prepare to write text in columns
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 10, ' ', tabwriter.AlignRight)

	// Find all wallpapers
	wallpapers, gnomeWallpapers, simpleTimedWallpapers := monitor.FindWallpapers()
	var collectionNames []string

	// Output all wallpaper collection names and paths (these are directories
	// with files of varying resolutions)
	for _, wp := range wallpapers {
		if wp.PartOfCollection {
			name := wp.CollectionName
			dir := filepath.Dir(wp.Path) + "/"
			if alsoPrintPath || !has(collectionNames, name) {
				fmt.Fprintf(w, "%s\t(%s)\t\t%s\n", name, "wallpaper collection", dir)
				collectionNames = append(collectionNames, wp.CollectionName)
			}
		}
	}

	// Timed wallpapers is a collection in the sense that it may point to
	// several wallpaper images.

	// Output all Simple Timed Wallpaper names and paths.
	for _, stw := range simpleTimedWallpapers {
		name := stw.Name
		path := stw.Path
		if alsoPrintPath || !has(collectionNames, name) {
			fmt.Fprintf(w, "%s\t(%s)\t\t%s\n", name, "simple timed wallpaper", path)
			collectionNames = append(collectionNames, name)
		}
	}

	// Output all GNOME timed wallpaper names and paths.
	for _, gw := range gnomeWallpapers {
		name := gw.CollectionName
		path := gw.Path
		if alsoPrintPath || !has(collectionNames, name) {
			fmt.Fprintf(w, "%s\t(%s)\t\t%s\n", name, "GNOME timed wallpaper", path)
			collectionNames = append(collectionNames, name)
		}
	}

	// Write the output to stdout
	w.Flush()
}
