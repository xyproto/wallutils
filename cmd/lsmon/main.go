package main

import (
	"fmt"
	"os"

	"github.com/xyproto/wallutils"
)

func main() {
	// Retrieve a slice of Monitor structs, or exit with an error
	monitors, err := wallutils.Monitors()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	// For every monitor, output the ID, width and height
	for _, mon := range monitors {
		if len(os.Args) > 1 && (os.Args[1] == "-dpi" || os.Args[1] == "-l") {
			fmt.Printf("%d: %dx%d (DPI: %dx%d)\n", mon.ID, mon.Width, mon.Height, mon.DPIw, mon.DPIh)
		} else if len(os.Args) > 1 && os.Args[1] == "--version" {
			fmt.Println(wallutils.VersionString)
			os.Exit(0)
		} else {
			fmt.Printf("%d: %dx%d\n", mon.ID, mon.Width, mon.Height)
		}
	}
}
