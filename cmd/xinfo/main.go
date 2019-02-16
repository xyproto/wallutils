package main

import (
	"fmt"
	"github.com/xyproto/wallutils"
	"os"
)

func main() {
	// Fetch the info string
	info, err := wallutils.XInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	// Output the info
	fmt.Println(info)
}
