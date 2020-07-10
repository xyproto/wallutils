package wallutils

import (
	"errors"
	"fmt"
)

// Mate windowmanager detector
type Mate struct {
	mode         string // none | wallpaper | centered | scaled | stretched | zoom | spanned, scaled is the default
	hasMate      bool
	hasGsettings bool
	hasChecked   bool
	verbose      bool
}

// Name returns the name of this window manager or desktop environment
func (m *Mate) Name() string {
	return "Mate"
}

// ExecutablesExists checks if executables associated with this backend exists in the PATH
func (m *Mate) ExecutablesExists() bool {
	// Cache the results
	m.hasGsettings = which("gsettings") != ""
	m.hasMate = which("mate-session") != ""
	m.hasChecked = true

	// The result may be used both outside of this file, and in SetWallpaper
	return m.hasMate && m.hasGsettings
}

// Running examines environment variables to try to figure out if this backend is currently running
func (m *Mate) Running() bool {
	return containsE("GDMSESSION", "mate") || containsE("XDG_SESSION_DESKTOP", "MATE") || containsE("XDG_CURRENT_DESKTOP", "MATE") || containsE("DESKTOP_SESSION", "xsessions/mate")
}

// SetMode will set the current way to display the wallpaper (stretched, tiled etc)
func (m *Mate) SetMode(mode string) {
	m.mode = mode
}

// SetVerbose can be used for setting the verbose field to true or false.
// This will cause this backend to output information about what is is doing on stdout.
func (m *Mate) SetVerbose(verbose bool) {
	m.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (m *Mate) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}
	// Check if dconf or gsettings are there, if we haven't already checked
	if !m.hasChecked {
		// This alters the state of m
		m.ExecutablesExists()
	}

	// initialize the mode setting (stretched/tiled etc)
	mode := defaultMode
	if m.mode != "" {
		mode = m.mode
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

	if !m.hasGsettings {
		return errors.New("could not find gsettings")
	}

	// Exit if the monitor configuration will cause artifacts when setting
	// the desktop wallpaper.
	NoXRandrOverlapOrExit(m.verbose)

	// Create a new GSettings struct, for dealing with GNOME settings
	g := NewGSettings("org.mate.background", m.verbose)

	// Set picture-options, if it is not already set to the desired value
	if g.Get("picture-options") != mode {
		if err := g.Set("picture-options", mode); err != nil {
			return err
		}
	}

	// Set the desktop wallpaper (also set it if it is already set, in case
	// the contents have changed)
	return g.Set("picture-filename", imageFilename)
}
