// +build cgo

package wallutils

import (
	"fmt"
)

// Sway windowmanager detector
type Sway struct {
	mode    string
	verbose bool
}

func (s *Sway) Name() string {
	return "Sway"
}

func (s *Sway) ExecutablesExists() bool {
	return which("sway") != "" && which("swaymsg") != ""
}

func (s *Sway) Running() bool {
	return hasE("SWAYSOCK") || (containsE("GDMSESSION", "sway") || containsE("XDG_SESSION_DESKTOP", "sway") || containsE("XDG_CURRENT_DESKTOP", "sway"))
}

func (s *Sway) SetMode(mode string) {
	s.mode = mode
}

func (s *Sway) SetVerbose(verbose bool) {
	s.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (s *Sway) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}
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
