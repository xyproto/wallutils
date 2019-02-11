package main

import (
	"fmt"
	"github.com/xyproto/monitor"
)

func main() {
	_, gnomeWallpapers, simpleTimedWallpapers := monitor.FindWallpapers()
	for _, stw := range simpleTimedWallpapers {
		fmt.Println("--- " + stw.Name + " ---")
		fmt.Println("path:", stw.Path)
		fmt.Println(stw)
		fmt.Println()
	}
	for _, gb := range gnomeWallpapers {
		fmt.Println(gb)
	}
}
