package main

import (
	"fmt"
	"github.com/xyproto/monitor"
	"os"
)

func main() {
	// Fetch the info string
	info, err := monitor.XInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	// Output the info
	fmt.Println(info)
}
