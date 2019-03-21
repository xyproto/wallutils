package wallutils

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"strconv"
	"strings"
)

// XrandrOverlap checks if the displays listed by "xrandr" overlaps, or not
// A slice of relevant lines from the xrandr output is also returned.
func XRandrOverlap() (bool, []string) {
	resolution_lines := []string{}
	if which("xrandr") == "" {
		return false, resolution_lines
	}
	xrandrOutput := output("xrandr", []string{}, false)
	rects := make([]*Rect, 0)
	for _, line := range strings.Split(xrandrOutput, "\n") {
		words := strings.Fields(line)
		if len(words) < 3 {
			continue
		}
		if words[1] != "connected" {
			continue
		}
		for _, res := range words[2:] {
			if !strings.Contains(res, "x") {
				continue
			}
			if strings.Count(res, "+") != 2 {
				continue
			}
			fields := strings.SplitN(res, "x", 2)
			ws, tail := fields[0], fields[1]
			fields = strings.SplitN(tail, "+", 2)
			hs, tail := fields[0], fields[1]
			fields = strings.SplitN(tail, "+", 2)
			xs, ys := fields[0], fields[1]

			resolution_lines = append(resolution_lines, line)
			//if verbose {
			//	fmt.Println("XRANDR: " + line)
			//}

			// Convert coordinates from string to int
			x, err := strconv.Atoi(xs)
			if err != nil {
				continue
			}
			y, err := strconv.Atoi(ys)
			if err != nil {
				continue
			}
			width, err := strconv.Atoi(ws)
			if err != nil {
				continue
			}
			height, err := strconv.Atoi(hs)
			if err != nil {
				continue
			}

			// Create a new Rect struct and append it to the collection
			r := NewRect(uint(x), uint(y), uint(width), uint(height))
			rects = append(rects, r)

			// Don't examine the rest of the words, but skip to the next line
			break
		}
	}
	// Check if the gathered display rectangles overlap
	return overlaps(rects), resolution_lines
}

// NoXrandrOverlapOrExit is a convenience function for making sure monitor
// configurations are not overlapping, as reported by "xrandr".
func NoXrandrOverlapOrExit(verbose bool) {
	if overlap, reslines := XRandrOverlap(); overlap {
		red := color.New(color.FgRed)
		white := color.New(color.FgWhite, color.Bold)
		red.Fprint(os.Stderr, "ERROR: ")
		fmt.Fprintln(os.Stderr, "xrandr shows overlapping monitor configurations:")
		white.Fprintln(os.Stderr, strings.Join(reslines, "\n"))
		os.Exit(1)
	} else if verbose {
		fmt.Println("No overlapping monitor configurations.")
	}
}
