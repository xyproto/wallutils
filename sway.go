//go:build cgo
// +build cgo

package wallutils

import (
	"fmt"

	"github.com/xyproto/env"
)

// Sway windowmanager detector
type Sway struct {
	mode    string
	verbose bool
}

// Name returns the name of this window manager or desktop environment
func (s *Sway) Name() string {
	return "Sway"
}

// ExecutablesExists checks if executables associated with this backend exists in the PATH
func (s *Sway) ExecutablesExists() bool {
	return which("sway") != "" && which("swaymsg") != ""
}

// Running examines environment variables to try to figure out if this backend is currently running
func (s *Sway) Running() bool {
	return env.Has("SWAYSOCK") || (env.Contains("GDMSESSION", "sway") || env.Contains("XDG_SESSION_DESKTOP", "sway") || env.Contains("XDG_CURRENT_DESKTOP", "sway"))
}

// SetMode will set the current way to display the wallpaper (stretched, tiled etc)
func (s *Sway) SetMode(mode string) {
	s.mode = mode
}

// SetVerbose can be used for setting the verbose field to true or false.
// This will cause this backend to output information about what is is doing on stdout.
func (s *Sway) SetVerbose(verbose bool) {
	s.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (s *Sway) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}

	// initialize the mode setting (stretched/tiled etc)
	mode := defaultMode
	if s.mode != "" {
		mode = s.mode
	}

	switch mode {
	case "center", "tile", "fill", "stretch":
		break
	case "scale", "scaled":
		mode = "fill"
	case "zoom", "zoomed", "stretched":
		mode = "stretch"
	default:
		// Invalid and unrecognized desktop wallpaper mode
		return fmt.Errorf("invalid desktop wallpaper mode for Sway: %s", mode)
	}

	return run("swaymsg", []string{"output * bg " + imageFilename + " " + mode}, s.verbose)
}
