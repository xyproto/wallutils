package main

import (
	"fmt"
	"github.com/xyproto/monitor"
	"os"
)

func main() {
	searchResults, err := monitor.FindWallpapers()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	for _, wp := range searchResults.Wallpapers() {
		fmt.Println(wp)
	}
}
