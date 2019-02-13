package monitor

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
	// TODO: To implement
	//return (containsE("GDMSESSION", "gnome2") || containsE("XDG_SESSION_DESKTOP", "gnome2") || containsE("XDG_CURRENT_DESKTOP", "gnome2"))
	return false
}

func (g2 *Gnome2) SetVerbose(verbose bool) {
	g2.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (g2 *Gnome2) SetWallpaper(imageFilename string) error {
	return run("gconftool-2 –type string –set /desktop/gnome/background/picture_filename \""+imageFilename+"\"", g2.verbose)
}
