package main

import (
	"fmt"
	"github.com/xyproto/monitor"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println(monitor.VersionString)
		os.Exit(0)
	}
	searchResults, err := monitor.FindWallpapers()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	withDetails := len(os.Args) > 1 && os.Args[1] == "-l"
	for _, wp := range searchResults.Wallpapers() {
		if withDetails {
			fmt.Println(wp)
		} else {
			fmt.Println(wp.Path)
		}
	}
}
