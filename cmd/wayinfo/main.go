package main

import (
	"fmt"
	"os"

	"github.com/xyproto/monitor"
)

func main() {
	// Fetch the info string
	info, err := monitor.WaylandInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	// Output the info
	fmt.Println(info)
}
