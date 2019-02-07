package monitor

// Gnome3 windowmanager detector
type Gnome3 struct {
}

func (s *Gnome3) Name() string {
	return "Gnome3"
}

func (s *Gnome3) ExecutablesExists() bool {
	return which("gsettings") != ""
}

func (s *Gnome3) Running() bool {
	// TODO: Needs testing
	return (containsE("GDMSESSION", "gnome") || containsE("XDG_SESSION_DESKTOP", "gnome") || containsE("XDG_CURRENT_DESKTOP", "gnome") || containsE("XDG_CURRENT_DESKTOP", "GNOME"))
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (s *Gnome3) SetWallpaper(imageFilename string) error {
	return run("gsettings set org.gnome.desktop.background picture-uri \"file://" + imageFilename + "\"")
}
