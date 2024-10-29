package wallutils

import (
	"errors"
	"fmt"

	"github.com/xyproto/env/v2"
)

// Deepin windowmanager detector
type Deepin struct {
	mode         string // none | wallpaper | centered | scaled | stretched | zoom | spanned, scaled is the default
	hasDeepin    bool
	hasGsettings bool
	hasChecked   bool
	verbose      bool
}

// Name returns the name of this window manager or desktop environment
func (d *Deepin) Name() string {
	return "Deepin"
}

// ExecutablesExists checks if executables associated with this backend exists in the PATH
func (d *Deepin) ExecutablesExists() bool {
	// Cache the results
	d.hasDeepin = which("deepin-session") != ""
	d.hasGsettings = which("gsettings") != ""
	d.hasChecked = true

	// The result may be used both outside of this file, and in SetWallpaper
	return d.hasDeepin && d.hasGsettings
}

// Running examines environment variables to try to figure out if this backend is currently running
func (d *Deepin) Running() bool {
	return env.Contains("GDMSESSION", "deepin") || env.Contains("XDG_CURRENT_DESKTOP", "Deepin") || env.Contains("DESKTOP_SESSION", "xsessions/deepin")
}

// SetMode will set the current way to display the wallpaper (stretched, tiled etc)
func (d *Deepin) SetMode(mode string) {
	d.mode = mode
}

// SetVerbose can be used for setting the verbose field to true or false.
// This will cause this backend to output information about what is is doing on stdout.
func (d *Deepin) SetVerbose(verbose bool) {
	d.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (d *Deepin) SetWallpaper(imageFilename string) error {
	// Check if the image exists
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}

	// Check if dconf or gsettings are there, if we haven't already checked
	if !d.hasChecked {
		// This alters the state of d
		d.ExecutablesExists()
	}

	// initialize the mode setting (stretched/tiled etc)
	mode := defaultMode
	if d.mode != "" {
		mode = d.mode
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
		return fmt.Errorf("invalid desktop wallpaper mode for Deepin: %s", mode)
	}

	if !d.hasGsettings {
		return errors.New("could not find gsettings")
	}

	// Exit if the monitor configuration will cause artifacts when setting
	// the desktop wallpaper.
	NoXRandrOverlapOrExit(d.verbose)

	// Create a new GSettings struct, for dealing with GNOME settings
	g := NewGSettings("com.deepin.wrap.gnome.desktop.background", d.verbose)

	// Set picture-options, if it is not already set to the desired value
	if g.Get("picture-options") != mode {
		if err := g.Set("picture-options", mode); err != nil {
			return err
		}
	}

	// Set the dark desktop wallpaper (also set it if it is already set)
	_ = g.Set("picture-uri-dark", "file://"+imageFilename)

	// Set the desktop wallpaper (also set it if it is already set)
	return g.Set("picture-uri", "file://"+imageFilename)
}
