package wallutils

import (
	"fmt"
)

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
	return (containsE("XDG_CURRENT_DESKTOP", "X-Cinnamon") || containsE("GDMSESSION", "cinnamon") || containsE("DESKTOP_SESSION", "xsessions/cinnamon"))
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
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}
	// Check if dconf or gsettings are there, if we haven't already checked
	if !c.hasChecked {
		c.ExecutablesExists()
	}

	mode := defaultMode

	// If c.mode is specified, do not use the default value
	if c.mode != "" {
		mode = c.mode
	}

	// possible values for gsettings / picture-options: "none", "wallpaper", "centered", "scaled", "stretched", "zoom", "spanned".
	switch mode {
	case "none", "wallpaper", "centered", "scaled", "stretched", "zoom", "spanned":
		break
	case "fill":
		// Invalid desktop wallpaper picture mode, use "stretched" instead
		mode = "stretched"
	default:
		// Invalid and unrecognized desktop wallpaper picture mode
		return fmt.Errorf("invalid desktop wallpaper picture mode for Cinnamon: %s", mode)
	}

	// Create a new GSettings struct, for dealing with GNOME settings
	g := NewGSettings("org.cinnamon.desktop.background", c.verbose)

	// Set picture-options, if it is not already set to the desired value
	if g.Get("picture-options") != mode {
		if err := g.Set("picture-options", mode); err != nil {
			return err
		}
	}

	// Set the desktop wallpaper (also set it if it is already set, in case
	// the contents have changed)
	return g.Set("picture-uri", "file://"+imageFilename)
}
