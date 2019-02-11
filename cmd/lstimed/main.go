package main

import (
	"fmt"
	"github.com/xyproto/monitor"
	"os"
)

func main() {
	alsoPrintPath := len(os.Args) > 1 && os.Args[1] == "-l"
	_, gnomeWallpapers := monitor.FindWallpapers()
	for _, gw := range gnomeWallpapers {
		if alsoPrintPath {
			numEvents := len(gw.Config.Statics) + len(gw.Config.Transitions)
			fmt.Printf("%s\t%s\tevents: %d\n", gw.CollectionName, gw.Path, numEvents)
		} else {
			fmt.Printf("%s\n", gw.CollectionName)
		}
	}
}
