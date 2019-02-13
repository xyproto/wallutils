package monitor

import (
	"fmt"
	"strings"
)

const simpleTimedWallpaperFormatVersion = "1.0"

// GnomeToSimple converts a Gnome Timed Wallpaper to a Simple Timed Wallpaper
func GnomeToSimple(gw *GnomeTimedWallpaper) (*SimpleTimedWallpaper, error) {

	// TODO: Convert from struct to struct, without excercising the serializer and the parser

	// Convert the given struct to the string contents of a SimpleTimedWallpaper file
	s, err := GnomeToSimpleString(gw)
	if err != nil {
		return nil, err
	}
	return DataToSimple(gw.Path, []byte(s))
}

// GnomeToSimpleString converts a Gnome Timed Wallpaper to a string
// representing a Simple Timed Wallpaper. The Path field in the given
// struct is not included in the output string.
func GnomeToSimpleString(gw *GnomeTimedWallpaper) (string, error) {
	//filename := gw.Path
	name := gw.Name

	var sb strings.Builder

	// Output the version of the format
	sb.WriteString("stw: " + simpleTimedWallpaperFormatVersion + "\n")

	// Output the name of the timed wallpaper
	sb.WriteString("name: " + name + "\n")

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
			return "", fmt.Errorf("element is not a <static> or <transition> tag: error: %s", err)
		}
		if s, ok := eInterface.(GStatic); ok {
			filenames = append(filenames, s.Filename)
		} else if t, ok := eInterface.(GTransition); ok {
			filenames = append(filenames, t.FromFilename)
			filenames = append(filenames, t.ToFilename)
		}
	}

	// Then find the common prefix and suffix of the image filenames
	commonPrefix := CommonPrefix(filenames)
	commonSuffix := CommonSuffix(filenames)

	// Output the format string
	sb.WriteString("format: " + commonPrefix + "%s" + commonSuffix + "\n")

	// Then output the timing information, for static images and for transitions
	for i := 0; i < totalElements; i++ {
		// The duration of the event is specified in the XML file, but not when it should start

		// Get an element, by index. This is an interface{} and is expected to be a GStatic or a GTransition
		eInterface, err := gw.Config.Get(i)
		if err != nil {
			return "", fmt.Errorf("element is not a <static> or <transition> tag: error: %s", err)
		}
		if s, ok := eInterface.(GStatic); ok {
			window := s.Duration()

			sb.WriteString(fmt.Sprintf("@%s: %s\n", c(eventTime), Meat(s.Filename, commonPrefix, commonSuffix)))

			// Increase the variable that keeps track of the time
			eventTime = eventTime.Add(window)

		} else if t, ok := eInterface.(GTransition); ok {
			window := t.Duration()
			from := eventTime
			upTo := eventTime.Add(window)

			if t.Type == "overlay" {
				sb.WriteString(fmt.Sprintf("@%s-%s: %s .. %s\n", c(from), c(upTo), Meat(t.FromFilename, commonPrefix, commonSuffix), Meat(t.ToFilename, commonPrefix, commonSuffix)))
			} else {
				sb.WriteString(fmt.Sprintf("@%s-%s: %s .. %s | %s\n", c(from), c(upTo), Meat(t.FromFilename, commonPrefix, commonSuffix), Meat(t.ToFilename, commonPrefix, commonSuffix), t.Type))
			}

			// Increase the variable that keeps track of the time
			eventTime = eventTime.Add(window)
		}
	}

	return strings.TrimSpace(sb.String()), nil
}

// GnomeFileToSimpleString reads and parses an XML file, then returns a string
// representing the contents of a Simple Timed Wallpaper file.
func GnomeFileToSimpleString(filename string) (string, error) {
	gw, err := ParseXML(filename)
	if err != nil {
		return "", fmt.Errorf("Could not parse %s: %s", filename, err)
	}
	return GnomeToSimpleString(gw)
}
