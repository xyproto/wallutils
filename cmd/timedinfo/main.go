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
	for _, stw := range searchResults.SimpleTimedWallpapers() {
		fmt.Println("--- " + stw.Name + " ---")
		fmt.Println("path:", stw.Path)
		fmt.Println(stw)
		fmt.Println()
	}
	for _, gb := range searchResults.GnomeTimedWallpapers() {
		fmt.Println("--- " + gb.Name + " ---")
		fmt.Println(gb)
	}
}
