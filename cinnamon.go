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

// sq single-quotes a string
func sq(s string) string {
	return "'" + s + "'"
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
	default:
		// Invalid desktop background mode, use "stretched"
		mode = "stretched"
	}

	// Set the desktop wallpaper picture mode
	if err := run("gsettings", []string{"set", "org.cinnamon.desktop.background", "picture-options", sq(mode)}, c.verbose); err != nil {
		return err
	}

	// Set the desktop wallpaper
	return run("gsettings", []string{"set", "org.cinnamon.desktop.background", "picture-uri", sq("file://" + imageFilename)}, c.verbose)
}
