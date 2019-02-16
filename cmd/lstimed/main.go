package main

import (
	"fmt"
	"github.com/xyproto/wallutils"
	"os"
	"text/tabwriter"
)

func main() {
	alsoPrintPath := len(os.Args) > 1 && os.Args[1] == "-l"

	// Prepare to write text in columns
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 10, ' ', tabwriter.AlignRight)

	searchResults, err := wallutils.FindWallpapers()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	for _, stw := range searchResults.SimpleTimedWallpapers() {
		if alsoPrintPath {
			numEvents := len(stw.Statics) + len(stw.Transitions)
			fmt.Fprintf(w, "%s\t%s\t\tevents: %d\n", stw.Name, stw.Path, numEvents)
		} else {
			fmt.Fprintf(w, "%s\n", stw.Name)
		}
	}
	for _, gw := range searchResults.GnomeTimedWallpapers() {
		if alsoPrintPath {
			numEvents := len(gw.Config.Statics) + len(gw.Config.Transitions)
			fmt.Fprintf(w, "%s\t%s\t\tevents: %d\n", gw.Name, gw.Path, numEvents)
		} else {
			fmt.Fprintf(w, "%s\n", gw.Name)
		}
	}
	w.Flush()
}
