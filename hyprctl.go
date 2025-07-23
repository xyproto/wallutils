package wallutils

import (
	"fmt"

	"github.com/xyproto/env/v2"
)

// Hyprctl compatible windowmanager
type Hyprctl struct {
	mode    string
	verbose bool
}

// Name returns the name of this window manager or desktop environment
func (h *Hyprctl) Name() string {
	return "Hyprctl"
}

// ExecutablesExists checks if executables associated with this backend exists in the PATH
func (h *Hyprctl) ExecutablesExists() bool {
	return which("hyprctl") != "" && which("hyprpaper") != ""
}

// Running examines environment variables to try to figure out if this backend is currently running.
func (h *Hyprctl) Running() bool {
	return env.Contains("XDG_SESSION_DESKTOP", "Hyprland") || env.Contains("DESKTOP_SESSION", "hyprland")
}

// SetMode will set the current way to display the wallpaper (stretched, tiled etc)
func (h *Hyprctl) SetMode(mode string) {
	h.mode = mode
}

// SetVerbose can be used for setting the verbose field to true or false.
// This will cause this backend to output information about what is is doing on stdout.
func (h *Hyprctl) SetVerbose(verbose bool) {
	h.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (h *Hyprctl) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}

	// Set the wallpaper mode
	mode := defaultMode
	if h.mode != "" {
		mode = h.mode
	}

	switch mode {
	case "stretched", "center", "fill", "fit", "scale", "scaled", "stretch", "tile", "zoom", "zoomed":
		mode = ""
	default:
		// Invalid and unrecognized desktop wallpaper mode
		return fmt.Errorf("invalid desktop wallpaper mode for swaybg: %s", mode)
	}

	// preload the wallpaper image using hyprctl
	err := run("hyprctl", []string{"hyprpaper", "preload", imageFilename}, h.verbose)
	if err != nil {
		return err
	}

	// reload the wallpaper image using hyprctl
	return run("hyprctl", []string{"hyprpaper", "reload", ",\"" + imageFilename + "\""}, h.verbose)
}
