package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xyproto/monitor"
)

const simpleTimedWallpaperFormatVersion = "1.0"

// c formats a timestamp as HH:MM
func c(t time.Time) string {
	return fmt.Sprintf("%.2d:%.2d", t.Hour(), t.Minute())
}

func CommonPrefix(sl []string) string {
	if len(sl) == 0 {
		return ""
	}
	shortestLength := len(sl[0])
	shortestString := sl[0]
	for _, s := range sl {
		if len(s) < shortestLength {
			shortestLength = len(s)
			shortestString = s
		}
	}
	if shortestLength == 0 {
		return ""
	}
	for i := 1; i < shortestLength; i++ {
		for _, s := range sl {
			if !strings.HasPrefix(s, shortestString[:i]) {
				return shortestString[:i-1]
			}
		}
	}
	return shortestString
}

func CommonSuffix(sl []string) string {
	if len(sl) == 0 {
		return ""
	}
	shortestLength := len(sl[0])
	shortestString := sl[0]
	for _, s := range sl {
		if len(s) < shortestLength {
			shortestLength = len(s)
			shortestString = s
		}
	}
	if shortestLength == 0 {
		return ""
	}
	for i := 1; i < shortestLength; i++ {
		for _, s := range sl {
			if !strings.HasSuffix(s, shortestString[shortestLength-i:]) {
				return shortestString[shortestLength-(i-1):]
			}
		}
	}
	return shortestString
}

// Meat returns the meat of the string: the part of the filename after the
// prefix and before the suffix
func Meat(filename, prefix, suffix string) string {
	if len(filename) < (len(prefix) + len(suffix)) {
		panic("too short filename")
	}
	return filename[len(prefix) : len(filename)-len(suffix)]
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "Please give the path to a GNOME timed wallpaper XML file as the first argument.")
		os.Exit(1)
	}

	filename := os.Args[1]

	gb, err := monitor.ParseXML(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse %s as XML: error: %s\n", filename, err)
		os.Exit(1)
	}

	// Output the version of the format
	fmt.Println("stw: " + simpleTimedWallpaperFormatVersion)

	// Use the name of the file, without the extension, as the name of this timed wallpaper
	name := filepath.Base(filename[:len(filename)-len(filepath.Ext(filename))])
	gw := monitor.NewGnomeWallpaper(name, filename, gb)

	// Output the name of the timed wallpaper
	fmt.Println("name: " + name)

	// Get the start time for the wallpaper collection (which is offset by X
	// seconds per static wallpaper)
	startTime := gw.StartTime()

	totalElements := len(gw.Config.Statics) + len(gw.Config.Transitions)

	// Keep track of the total time. It is increased every time a new element duration is encountered.
	eventTime := startTime

	// First, only gather all the image filenames
	var filenames []string
	for i := 0; i < totalElements; i++ {
		// Get an element, by index. This is an interface{} and is expected to be a GStatic or a GTransition
		eInterface, err := gw.Config.Get(i)
		if err != nil {
			fmt.Fprintf(os.Stderr, "element is not a <static> or <transition> tag: error: %s\n", err)
			os.Exit(1)
		}
		if s, ok := eInterface.(monitor.GStatic); ok {
			filenames = append(filenames, s.Filename)
		} else if t, ok := eInterface.(monitor.GTransition); ok {
			filenames = append(filenames, t.FromFilename)
			filenames = append(filenames, t.ToFilename)
		}
	}

	// Then find the common prefix and suffix of the image filenames
	commonPrefix := CommonPrefix(filenames)
	commonSuffix := CommonSuffix(filenames)

	// Output the format string
	fmt.Println("format: " + commonPrefix + "%s" + commonSuffix)

	// Then output the timing information, for static images and for transitions
	for i := 0; i < totalElements; i++ {
		// The duration of the event is specified in the XML file, but not when it should start

		// Get an element, by index. This is an interface{} and is expected to be a GStatic or a GTransition
		eInterface, err := gw.Config.Get(i)
		if err != nil {
			fmt.Fprintf(os.Stderr, "element is not a <static> or <transition> tag: error: %s\n", err)
			os.Exit(1)
		}
		if s, ok := eInterface.(monitor.GStatic); ok {
			window := s.Duration()

			fmt.Printf("@%s: %s\n", c(eventTime), Meat(s.Filename, commonPrefix, commonSuffix))

			// Increase the variable that keeps track of the time
			eventTime = eventTime.Add(window)

		} else if t, ok := eInterface.(monitor.GTransition); ok {
			window := t.Duration()
			from := eventTime
			upTo := eventTime.Add(window)

			if t.Type == "overlay" {
				fmt.Printf("@%s-%s: %s .. %s\n", c(from), c(upTo), Meat(t.FromFilename, commonPrefix, commonSuffix), Meat(t.ToFilename, commonPrefix, commonSuffix))
			} else {
				fmt.Printf("@%s-%s: %s .. %s | %s\n", c(from), c(upTo), Meat(t.FromFilename, commonPrefix, commonSuffix), Meat(t.ToFilename, commonPrefix, commonSuffix), t.Type)
			}

			// Increase the variable that keeps track of the time
			eventTime = eventTime.Add(window)
		}
	}

}
