package wallutils

import (
	"errors"
	"fmt"
)

// Cinnamon windowmanager detector
type Cinnamon struct {
	mode         string // none | wallpaper | centered | scaled | stretched | zoom | spanned, scaled is the default
	hasCinnamon  bool
	hasGsettings bool
	hasChecked   bool
	verbose      bool
}

// Name returns the name of this window manager or desktop environment
func (c *Cinnamon) Name() string {
	return "Cinnamon"
}

// ExecutablesExists checks if executables associated with this backend exists in the PATH
func (c *Cinnamon) ExecutablesExists() bool {
	// Cache the results
	c.hasGsettings = which("gsettings") != ""
	c.hasCinnamon = which("cinnamon") != ""
	c.hasChecked = true

	// The result may be used both outside of this file, and in SetWallpaper
	return c.hasCinnamon && c.hasGsettings
}

// Running examines environment variables to try to figure out if this backend is currently running
func (c *Cinnamon) Running() bool {
	return (containsE("XDG_CURRENT_DESKTOP", "X-Cinnamon") || containsE("GDMSESSION", "cinnamon") || containsE("DESKTOP_SESSION", "xsessions/cinnamon"))
}

// SetMode will set the current way to display the wallpaper (stretched, tiled etc)
func (c *Cinnamon) SetMode(mode string) {
	c.mode = mode
}

// SetVerbose can be used for setting the verbose field to true or false.
// This will cause this backend to output information about what is is doing on stdout.
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
		// This alters the state of c
		c.ExecutablesExists()
	}

	// initialize the mode setting (stretched/tiled etc)
	mode := defaultMode
	if c.mode != "" {
		mode = c.mode
	}

	// possible values for gsettings / picture-options: "none", "wallpaper", "centered", "scaled", "stretched", "zoom", "spanned".
	switch mode {
	case "none", "wallpaper", "centered", "scaled", "stretched", "zoom", "spanned":
		break
	case "stretch":
		mode = "stretched"
	case "center":
		mode = "centered"
	case "fill", "scale":
		mode = "scaled"
	case "tile":
		mode = "wallpaper"
	default:
		// Invalid and unrecognized desktop wallpaper mode
		return fmt.Errorf("invalid desktop wallpaper mode for MATE: %s", mode)
	}

	if !c.hasGsettings {
		return errors.New("could not find gsettings")
	}

	// Exit if the monitor configuration will cause artifacts when setting
	// the desktop wallpaper.
	NoXRandrOverlapOrExit(c.verbose)

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
