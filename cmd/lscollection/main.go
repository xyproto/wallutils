package main

import (
	"fmt"
	"github.com/xyproto/wallutils"
	"os"
	"path/filepath"
	"text/tabwriter"
)

// has checks if the given string slice contains the given string
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

	// Find all wallpapers
	searchResults, err := wallutils.FindWallpapers()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if !alsoPrintPath {
		for _, name := range searchResults.CollectionNames() {
			fmt.Println(name)
		}
		return
	}

	// Prepare to write text in columns
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 10, ' ', tabwriter.AlignRight)

	// Output all wallpaper collection names and paths (these are directories
	// with files of varying resolutions)
	var collectionNames []string
	for _, wp := range searchResults.Wallpapers() {
		if wp.PartOfCollection {
			name := wp.CollectionName
			dir := filepath.Dir(wp.Path) + "/"
			if alsoPrintPath && !has(collectionNames, name) {
				fmt.Fprintf(w, "%s\t%s\t\t%s\n", name, "Wallpaper Collection", dir)
				collectionNames = append(collectionNames, wp.CollectionName)
			}
		}
	}

	// Timed wallpapers is a collection in the sense that it may point to
	// several wallpaper images.

	// Output all Simple Timed Wallpaper names and paths.
	collectionNames = []string{}
	for _, stw := range searchResults.SimpleTimedWallpapers() {
		name := stw.Name
		path := stw.Path
		if alsoPrintPath && !has(collectionNames, name) {
			fmt.Fprintf(w, "%s\t%s\t\t%s\n", name, "Simple Timed Wallpaper", path)
			collectionNames = append(collectionNames, name)
		}
	}

	// Output all GNOME timed wallpaper names and paths.
	collectionNames = []string{}
	for _, gw := range searchResults.GnomeTimedWallpapers() {
		name := gw.Name
		path := gw.Path
		if alsoPrintPath && !has(collectionNames, name) {
			fmt.Fprintf(w, "%s\t%s\t\t%s\n", name, "GNOME Timed Wallpaper", path)
			collectionNames = append(collectionNames, name)
		}
	}

	// Write the output to stdout
	w.Flush()
}
