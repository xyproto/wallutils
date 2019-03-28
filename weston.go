package wallutils

import (
	"errors"
	"fmt"
	"os"
)

// Weston windowmanager detector
type Weston struct {
	mode    string
	verbose bool
}

func (w *Weston) Name() string {
	return "Weston"
}

func (w *Weston) ExecutablesExists() bool {
	return which("weston") != ""
}

func (w *Weston) Running() bool {
	return containsE("GDMSESSION", "weston") || containsE("XDG_SESSION_DESKTOP", "weston") || containsE("XDG_CURRENT_DESKTOP", "weston")
}

func (w *Weston) SetMode(mode string) {
	w.mode = mode
}

func (w *Weston) SetVerbose(verbose bool) {
	w.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (*Weston) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}

	fmt.Println("WESTON CONFIG FILE: ", os.Getenv("WESTON_CONFIG_FILE"))

	// TODO: Add the following to ~/.config/weston.ini
	//       (or whichever configuration file Weston uses)
	// [shell]
	// background-image=/home/user/somewhere/image.jpg
	// background-type=scale
	//
	// Also use w.mode for setting the background type

	return errors.New("Weston currently does not support changing the desktop wallpaper at runtime")
}
