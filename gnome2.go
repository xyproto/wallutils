package monitor

// Gnome2 windowmanager detector
type Gnome2 struct {
}

func (s *Gnome2) Name() string {
	return "Gnome2"
}

func (s *Gnome2) ExecutablesExists() bool {
	return which("gconftool-2") != ""
}

func (s *Gnome2) Running() bool {
	// TODO: To implement
	//return (containsE("GDMSESSION", "gnome2") || containsE("XDG_SESSION_DESKTOP", "gnome2") || containsE("XDG_CURRENT_DESKTOP", "gnome2"))
	return false
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (s *Gnome2) SetWallpaper(imageFilename string) error {
	return run("gconftool-2 –type string –set /desktop/gnome/background/picture_filename \"" + imageFilename + "\"")
}
