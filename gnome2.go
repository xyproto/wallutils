package wallutils

import (
	"fmt"
	"os"
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
	return (os.Getenv("GDMSESSION") == "gnome") || (os.Getenv("DESKTOP_SESSION") == "gnome")
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
	// TODO: Confirm that this works and find a way to set the mode (like "tile" or "fill")
	return run("gconftool-2", []string{"--type", "string", "--set", "/desktop/gnome/background/picture_filename", imageFilename}, g2.verbose)
}
