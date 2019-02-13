package monitor

// Mate windowmanager detector
type Mate struct {
	mode         string // none | wallpaper | centered | scaled | stretched | zoom | spanned, scaled is the default
	hasDconf     bool
	hasGsettings bool
	hasChecked   bool
	verbose      bool
}

func (m *Mate) Name() string {
	return "Mate"
}

func (m *Mate) ExecutablesExists() bool {
	m.hasDconf = which("dconf") != ""
	m.hasGsettings = which("gsettings") != ""
	m.hasChecked = true
	return m.hasDconf || m.hasGsettings
}

func (m *Mate) Running() bool {
	return (containsE("GDMSESSION", "mate") || containsE("XDG_SESSION_DESKTOP", "mate") || containsE("XDG_CURRENT_DESKTOP", "mate"))
}

func (m *Mate) SetMode(mode string) {
	m.mode = mode
}

func (m *Mate) SetVerbose(verbose bool) {
	m.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (m *Mate) SetWallpaper(imageFilename string) error {
	// Check if dconf or gsettings are there, if we haven't already checked
	if !m.hasChecked {
		m.ExecutablesExists()
	}
	// Set the desktop wallpaper picture mode
	mode := defaultMode
	if m.mode != "" {
		mode = m.mode
	}
	// Change the background with either dconf or gsettings
	if m.hasDconf {
		// use dconf
		if err := run("dconf write /org/mate/desktop/background/picture-filename \"'"+imageFilename+"'\"", m.verbose); err != nil {
			return err
		}
		return run("dconf write /org/mate/desktop/background/picture-options \"'"+mode+"'\"", m.verbose)
	}
	// use gsettings
	if err := run("gsettings set org.mate.background picture-filename '"+imageFilename+"'", m.verbose); err != nil {
		return err
	}
	return run("gsettings set org.mate.background picture-options '"+mode+"'", m.verbose)

}
