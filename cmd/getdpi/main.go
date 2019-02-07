package main

import (
	"fmt"
	"os"

	"github.com/xyproto/monitor"
)

const versionString = "1.2.0"

func main() {
	// Retrieve a slice of Monitor structs, or exit with an error
	monitors, err := monitor.Detect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	// Output the average DPI
	DPIw, DPIh := uint(0), uint(0)
	for _, monitor := range monitors {
		DPIw += monitor.DPIw
		DPIh += monitor.DPIh
	}
	DPIw /= uint(len(monitors))
	DPIh /= uint(len(monitors))

	// Check if -b is given (for outputting both numbers)
	if len(os.Args) > 1 && os.Args[1] == "-b" {
		fmt.Printf("%dx%d\n", DPIw, DPIh)
		return
	} else if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println(versionString)
		return
	}

	// Output a single number
	fmt.Println(DPIw)
}
