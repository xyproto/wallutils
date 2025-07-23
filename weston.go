package wallutils

import (
	"errors"
	"fmt"

	"github.com/xyproto/env/v2"
)

// Weston windowmanager detector
type Weston struct {
	mode    string
	verbose bool
}

// Name returns the name of this window manager or desktop environment
func (w *Weston) Name() string {
	return "Weston"
}

// ExecutablesExists checks if executables associated with this backend exists in the PATH
func (w *Weston) ExecutablesExists() bool {
	return which("weston") != ""
}

// Running examines environment variables to try to figure out if this backend is currently running
func (w *Weston) Running() bool {
	return env.Contains("GDMSESSION", "weston") || env.Contains("XDG_SESSION_DESKTOP", "weston") || env.Contains("XDG_CURRENT_DESKTOP", "weston")
}

// SetMode will set the current way to display the wallpaper (stretched, tiled etc)
func (w *Weston) SetMode(mode string) {
	w.mode = mode
}

// SetVerbose can be used for setting the verbose field to true or false.
// This will cause this backend to output information about what is is doing on stdout.
func (w *Weston) SetVerbose(verbose bool) {
	w.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (*Weston) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}

	fmt.Println("WESTON CONFIG FILE: ", env.Str("WESTON_CONFIG_FILE"))

	// TODO: Add the following to ~/.config/weston.ini
	//       (or whichever configuration file Weston uses)
	// [shell]
	// background-image=/home/user/somewhere/image.jpg
	// background-type=scale
	//
	// Also use w.mode for setting the background type

	return errors.New("Weston currently does not support changing the desktop wallpaper at runtime")
}
