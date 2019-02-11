package main

import (
	"fmt"
	"github.com/xyproto/monitor"
	"os"
	"text/tabwriter"
)

func main() {
	alsoPrintPath := len(os.Args) > 1 && os.Args[1] == "-l"

	// Prepare to write text in columns
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 10, ' ', tabwriter.AlignRight)

	_, gnomeWallpapers, simpleTimedWallpapers := monitor.FindWallpapers()
	for _, stw := range simpleTimedWallpapers {
		if alsoPrintPath {
			numEvents := len(stw.Statics) + len(stw.Transitions)
			fmt.Fprintf(w, "%s\t%s\t\tevents: %d\n", stw.Name, stw.Path, numEvents)
		} else {
			fmt.Fprintf(w, "%s\n", stw.Name)
		}
	}
	for _, gw := range gnomeWallpapers {
		if alsoPrintPath {
			numEvents := len(gw.Config.Statics) + len(gw.Config.Transitions)
			fmt.Fprintf(w, "%s\t%s\t\tevents: %d\n", gw.CollectionName, gw.Path, numEvents)
		} else {
			fmt.Fprintf(w, "%s\n", gw.CollectionName)
		}
	}

	w.Flush()
}
