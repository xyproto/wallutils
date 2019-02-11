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
	wallpapers, gnomeWallpapers := monitor.FindWallpapers()
	var collectionNames []string

	// Output all wallpaper collection names and paths (these are directories
	// with files of varying resolutions)
	for _, wp := range wallpapers {
		if wp.PartOfCollection {
			name := wp.CollectionName
			dir := filepath.Dir(wp.Path) + "/"
			if !has(collectionNames, name) {
				fmt.Fprintf(w, "%s\t(%s)\t\t%s\n", name, "wallpaper collection", dir)
				collectionNames = append(collectionNames, wp.CollectionName)
			}
		}
	}

	// Output all timed wallpaper names and paths.
	// Timed wallpapers is a collection in the sense that it may point to
	// several wallpaper images.
	for _, gw := range gnomeWallpapers {
		name := gw.CollectionName
		path := gw.Path
		if !has(collectionNames, name) {
			fmt.Fprintf(w, "%s\t(%s)\t\t%s\n", name, "timed wallpaper", path)
			collectionNames = append(collectionNames, name)
		}
	}

	// Write the output to stdout
	w.Flush()
}
