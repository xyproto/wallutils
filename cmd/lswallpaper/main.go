package main

import (
	"fmt"
	"github.com/xyproto/monitor"
)

func main() {
	wallpapers, _ := monitor.FindWallpapers()
	for _, wp := range wallpapers {
		fmt.Println(wp)
	}
}
