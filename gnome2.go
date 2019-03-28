package wallutils

import (
	"fmt"
)

// Gnome2 windowmanager detector
type Gnome2 struct {
	mode    string
	verbose bool
}

func (g2 *Gnome2) Name() string {
	return "Gnome2"
}

func (g2 *Gnome2) ExecutablesExists() bool {
	return which("gconftool-2") != ""
}

func (g2 *Gnome2) Running() bool {
	return (containsE("GDMSESSION", "gnome2") || containsE("XDG_SESSION_DESKTOP", "gnome2") || containsE("XDG_CURRENT_DESKTOP", "gnome2"))
}

func (g2 *Gnome2) SetMode(mode string) {
	g2.mode = mode
}

func (g2 *Gnome2) SetVerbose(verbose bool) {
	g2.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (g2 *Gnome2) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}
	// TODO: Confirm that this works
	// TODO: Find out how to set the wallpaper mode as well
	return run("gconftool-2", []string{"–type", "string", "–set", "/desktop/gnome/background/picture_filename", imageFilename}, g2.verbose)
}
