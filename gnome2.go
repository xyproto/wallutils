package wallutils

import (
	"fmt"
)

// TODO: Confirm that this is working under GNOME 2

// Gnome2 windowmanager detector
type Gnome2 struct {
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

func (g2 *Gnome2) SetVerbose(verbose bool) {
	g2.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (g2 *Gnome2) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}
	return run("gconftool-2", []string{"–type", "string", "–set", "/desktop/gnome/background/picture_filename", imageFilename}, g2.verbose)
}
