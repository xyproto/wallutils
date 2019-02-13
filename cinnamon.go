package monitor

// Cinnamon windowmanager detector
type Cinnamon struct {
	mode         string // none | wallpaper | centered | scaled | stretched | zoom | spanned, scaled is the default
	hasDconf     bool
	hasGsettings bool
	hasChecked   bool
	verbose      bool
}

func (c *Cinnamon) Name() string {
	return "Cinnamon"
}

func (c *Cinnamon) ExecutablesExists() bool {
	c.hasDconf = which("dconf") != ""
	c.hasGsettings = which("gsettings") != ""
	c.hasChecked = true
	return c.hasDconf || c.hasGsettings
}

func (c *Cinnamon) Running() bool {
	return (containsE("GDMSESSION", "cinnamon") || containsE("XDG_SESSION_DESKTOP", "cinnamon") || containsE("XDG_CURRENT_DESKTOP", "cinnamon"))
}

func (c *Cinnamon) SetMode(mode string) {
	c.mode = mode
}

func (c *Cinnamon) SetVerbose(verbose bool) {
	c.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (c *Cinnamon) SetWallpaper(imageFilename string) error {
	// Check if dconf or gsettings are there, if we haven't already checked
	if !c.hasChecked {
		c.ExecutablesExists()
	}
	// Set the desktop wallpaper picture mode
	mode := defaultMode
	if c.mode != "" {
		mode = c.mode
	}
	// Change the background with either dconf or gsettings
	if c.hasDconf {
		// use dconf
		if err := run("dconf write /org/cinnamon/desktop/background/picture-filename \"'"+imageFilename+"'\"", c.verbose); err != nil {
			return err
		}
		return run("dconf write /org/cinnamon/desktop/background/picture-options \"'"+mode+"'\"", c.verbose)
	}
	// use gsettings
	if err := run("gsettings set org.cinnamon.background picture-filename '"+imageFilename+"'", c.verbose); err != nil {
		return err
	}
	return run("gsettings set org.cinnamon.background picture-options '"+mode+"'", c.verbose)
}
