// Package wallutils can deal with monitors, resolution, dpi, wallpapers, wallpaper collections, timed wallpapers and converting GNOME timed wallpapers to the Simple Timed Wallpaper format.
package wallutils

import (
	"fmt"
	"os"
	"strings"
)

// VersionString is the current version of wallutils and all included utilities
const VersionString = "5.3.0"

// Quit with a nicely formatted error message to stderr
func Quit(err error) {
	msg := err.Error()
	if !strings.HasSuffix(msg, ".") && !strings.HasSuffix(msg, "!") && !strings.Contains(msg, ":") {
		msg += "."
	}
	fmt.Fprintf(os.Stderr, "%s%s\n", strings.ToUpper(string(msg[0])), msg[1:])
	os.Exit(1)
}
