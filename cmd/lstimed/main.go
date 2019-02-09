package main

import (
	"fmt"
	"github.com/xyproto/monitor"
)

func main() {
	_, gnomeWallpapers := monitor.FindWallpapers()
	for _, gb := range gnomeWallpapers {
		fmt.Println(gb.CollectionName)
	}
}
