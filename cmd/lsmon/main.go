package main

import (
	"fmt"
	"os"

	"github.com/xyproto/monitor"
)

func main() {
	// Retrieve a slice of Monitor structs, or exit with an error
	monitors, err := monitor.Detect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	// For every monitor, output the ID, width and height
	for _, monitor := range monitors {
		if len(os.Args) > 1 && os.Args[1] == "-dpi" {
			fmt.Printf("%d: %dx%d (DPI: %dx%d)\n", monitor.ID, monitor.Width, monitor.Height, monitor.DPIw, monitor.DPIh)
		} else {
			fmt.Printf("%d: %dx%d\n", monitor.ID, monitor.Width, monitor.Height)
		}
	}
}
