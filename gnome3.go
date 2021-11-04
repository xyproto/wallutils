package wallutils

import (
	"errors"
	"fmt"

	"github.com/xyproto/env"
)

// Gnome3 windowmanager detector
type Gnome3 struct {
	mode         string // none | wallpaper | centered | scaled | stretched | zoom | spanned, scaled is the default
	hasGnome3    bool
	hasGsettings bool
	hasChecked   bool
	verbose      bool
}

// Name returns the name of this window manager or desktop environment
func (g3 *Gnome3) Name() string {
	return "Gnome3"
}

// ExecutablesExists checks if executables associated with this backend exists in the PATH
func (g3 *Gnome3) ExecutablesExists() bool {
	// Cache the results
	g3.hasGsettings = which("gsettings") != ""
	g3.hasGnome3 = which("gnome-session") != ""
	g3.hasChecked = true

	// The result may be used both outside of this file, and in SetWallpaper
	return g3.hasGnome3 && g3.hasGsettings
}

// Running examines environment variables to try to figure out if this backend is currently running
func (g3 *Gnome3) Running() bool {
	return (env.Contains("GDMSESSION", "gnome") || env.Contains("XDG_SESSION_DESKTOP", "gnome"))
}

// SetMode will set the current way to display the wallpaper (stretched, tiled etc)
func (g3 *Gnome3) SetMode(mode string) {
	g3.mode = mode
}

// SetVerbose can be used for setting the verbose field to true or false.
// This will cause this backend to output information about what is is doing on stdout.
func (g3 *Gnome3) SetVerbose(verbose bool) {
	g3.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (g3 *Gnome3) SetWallpaper(imageFilename string) error {
	// Check if the image exists
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}

	// Check if dconf or gsettings are there, if we haven't already checked
	if !g3.hasChecked {
		// This alters the state of g3
		g3.ExecutablesExists()
	}

	// initialize the mode setting (stretched/tiled etc)
	mode := defaultMode
	if g3.mode != "" {
		mode = g3.mode
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
		return fmt.Errorf("invalid desktop wallpaper mode for GNOME3: %s", mode)
	}

	if !g3.hasGsettings {
		return errors.New("could not find gsettings")
	}

	// Exit if the monitor configuration will cause artifacts when setting
	// the desktop wallpaper.
	NoXRandrOverlapOrExit(g3.verbose)

	// Create a new GSettings struct, for dealing with GNOME settings
	g := NewGSettings("org.gnome.desktop.background", g3.verbose)

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
