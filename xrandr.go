package wallutils

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"os"
	"strconv"
	"strings"
)

type XRandr struct {
	hasOverlap      bool
	resolutionLines []string
	verbose         bool
	hasChecked      bool
}

func NewXRandr(verbose bool) (*XRandr, error) {
	if which("xrandr") == "" {
		return nil, errors.New("could not find the xrandr executable")
	}
	x := &XRandr{verbose: verbose}
	x.CheckOverlap()
	return x, nil
}

// Reset the XRander check, for preparing to run xrandr again
func (x *XRandr) Reset() {
	x.hasChecked = false
	x.resolutionLines = []string{}
}

// CheckOverlap checks if the displays listed by "xrandr" overlaps, or not
// A slice of relevant lines from the xrandr output is stored in the struct.
func (x *XRandr) CheckOverlap() {
	if x.hasChecked {
		return
	}
	if x.verbose {
		fmt.Print("Running ")
	}
	xrandrOutput := output("xrandr", []string{}, x.verbose)
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

			x.resolutionLines = append(x.resolutionLines, line)

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
	x.hasOverlap = overlaps(rects)
	x.hasChecked = true
}

// String returns a multiline string with the collected
// resolution lines from xrandr (if any).
func (x *XRandr) String() string {
	return strings.Join(x.resolutionLines, "\n")
}

// QuitIfOverlap will quit with an error if monitor configurations overlap
func (x *XRandr) QuitIfOverlap() {
	if x.hasOverlap {
		red := color.New(color.FgRed)
		white := color.New(color.FgWhite, color.Bold)
		red.Fprint(os.Stderr, "ERROR: ")
		fmt.Fprintln(os.Stderr, "xrandr shows overlapping monitor configurations:")
		white.Fprintln(os.Stderr, x)
		os.Exit(1)
	}
}

var cachedXRandr *XRandr

// NoXRandrOverlapOrExit is a convenience function for making sure monitor
// configurations are not overlapping, as reported by "xrandr".
func NoXRandrOverlapOrExit(verbose bool) {
	var (
		err        error
		initialRun bool
	)
	if cachedXRandr == nil {
		cachedXRandr, err = NewXRandr(verbose)
		if err != nil {
			// Could not check, just return
			return
		}
		initialRun = true
	}

	// Exit with an error if monitor configurations are overlapping
	cachedXRandr.QuitIfOverlap()

	if initialRun && cachedXRandr.verbose {
		fmt.Println("Detected no overlapping monitor configurations.")
	}
}
