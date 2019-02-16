package main

import (
	"fmt"
	"github.com/xyproto/wallutils"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println(wallutils.VersionString)
		os.Exit(0)
	}
	searchResults, err := wallutils.FindWallpapers()
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
