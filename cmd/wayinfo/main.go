package main

import (
	"fmt"
	"os"

	"github.com/xyproto/wallutils"
)

func main() {
	// Fetch the info string
	info, err := wallutils.WaylandInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	// Output the info
	fmt.Println(info)
}
