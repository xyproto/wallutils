package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/xyproto/monitor"
	"strings"

	"os"
)

// Indent all lines with the given prefix.
// Will trim the right side of the string for newlines before indenting.
func Indent(s string, prefix string) string {
	return prefix + strings.Replace(strings.TrimRight(s, "\n"), "\n", "\n"+prefix, -1)
}

func main() {
	searchResults, err := monitor.FindWallpapers()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	white := color.New(color.FgWhite, color.Bold)
	blue := color.New(color.FgBlue, color.Bold)
	gray := color.New(color.Reset)
	const prefix = "\t"
	first := true
	for _, stw := range searchResults.SimpleTimedWallpapers() {
		if first {
			first = false
		} else {
			fmt.Println()
		}
		white.Print("Simple Timed Wallpaper: ")
		blue.Print(stw.Name)
		fmt.Println()
		gray.Println("\n" + Indent("path: "+stw.Path+"\n"+stw.String(), prefix))
	}
	for _, gtw := range searchResults.GnomeTimedWallpapers() {
		if first {
			first = false
		} else {
			fmt.Println()
		}
		white.Print("GNOME Timed Wallpaper: ")
		blue.Print(gtw.Name)
		fmt.Println()
		gray.Println("\n" + Indent(gtw.String(), prefix))
	}
}
