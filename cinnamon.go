package monitor

// Cinnamon windowmanager detector
type Cinnamon struct {
	mode         string // none | wallpaper | centered | scaled | stretched | zoom | spanned, scaled is the default
	hasDconf     bool
	hasGsettings bool
	hasChecked   bool
}

func (s *Cinnamon) Name() string {
	return "Cinnamon"
}

func (s *Cinnamon) ExecutablesExists() bool {
	s.hasDconf = which("dconf") != ""
	s.hasGsettings = which("gsettings") != ""
	s.hasChecked = true
	return s.hasDconf || s.hasGsettings
}

func (s *Cinnamon) Running() bool {
	// TODO: To test
	return (containsE("GDMSESSION", "cinnamon") || containsE("XDG_SESSION_DESKTOP", "cinnamon") || containsE("XDG_CURRENT_DESKTOP", "cinnamon"))
}

func (s *Cinnamon) SetMode(mode string) {
	s.mode = mode
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (s *Cinnamon) SetWallpaper(imageFilename string) error {
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
		if err := run("dconf write /org/cinnamon/desktop/background/picture-filename \"'" + imageFilename + "'\""); err != nil {
			return err
		}
		return run("dconf write /org/cinnamon/desktop/background/picture-options \"'" + mode + "'\"")
	}
	// use gsettings
	if err := run("gsettings set org.cinnamon.background picture-filename '" + imageFilename + "'"); err != nil {
		return err
	}
	return run("gsettings set org.cinnamon.background picture-options '" + mode + "'")
}
