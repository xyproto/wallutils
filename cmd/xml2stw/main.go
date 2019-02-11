package main

import (
	"fmt"
	"os"

	"github.com/xyproto/monitor"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "Please give the path to a GNOME timed wallpaper XML file as the first argument.")
		os.Exit(1)
	}
	filename := os.Args[1]

	s, err := monitor.GnomeXMLToSimpleTimed(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(s)
}
