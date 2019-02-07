package monitor

// Mate windowmanager detector
type Mate struct {
	mode         string // none | wallpaper | centered | scaled | stretched | zoom | spanned, scaled is the default
	hasDconf     bool
	hasGsettings bool
	hasChecked   bool
}

func (s *Mate) Name() string {
	return "Mate"
}

func (s *Mate) ExecutablesExists() bool {
	s.hasDconf = which("dconf") != ""
	s.hasGsettings = which("gsettings") != ""
	s.hasChecked = true
	return s.hasDconf || s.hasGsettings
}

func (s *Mate) Running() bool {
	return (containsE("GDMSESSION", "mate") || containsE("XDG_SESSION_DESKTOP", "mate") || containsE("XDG_CURRENT_DESKTOP", "mate"))
}

func (s *Mate) SetMode(mode string) {
	s.mode = mode
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (s *Mate) SetWallpaper(imageFilename string) error {
	// Check if dconf or gsettings are there, if we haven't already checked
	if !s.hasChecked {
		s.ExecutablesExists()
	}
	// Set the desktop wallpaper picture mode
	mode := defaultMode
	if s.mode != "" {
		mode = s.mode
	}
	// Change the background with either dconf or gsettings
	if s.hasDconf {
		// use dconf
		if err := run("dconf write /org/mate/desktop/background/picture-filename \"'" + imageFilename + "'\""); err != nil {
			return err
		}
		return run("dconf write /org/mate/desktop/background/picture-options \"'" + mode + "'\"")
	}
	// use gsettings
	if err := run("gsettings set org.mate.background picture-filename '" + imageFilename + "'"); err != nil {
		return err
	}
	return run("gsettings set org.mate.background picture-options '" + mode + "'")

}
