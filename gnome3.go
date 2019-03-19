package wallutils

import (
	"errors"
)

// Gnome3 windowmanager detector
type Gnome3 struct {
	verbose bool
}

func (g3 *Gnome3) Name() string {
	return "Gnome3"
}

func (g3 *Gnome3) ExecutablesExists() bool {
	return which("gsettings") != ""
}

func (g3 *Gnome3) Running() bool {
	return (containsE("GDMSESSION", "gnome") || containsE("XDG_SESSION_DESKTOP", "gnome") || containsE("XDG_CURRENT_DESKTOP", "gnome") || containsE("XDG_CURRENT_DESKTOP", "GNOME"))
}

func (g3 *Gnome3) SetVerbose(verbose bool) {
	g3.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (g3 *Gnome3) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return errors.New(imageFilename + " does not exist")
	}
	return runShell("gsettings set org.gnome.desktop.background picture-uri \"file://"+imageFilename+"\"", g3.verbose)
}
