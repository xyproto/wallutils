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

func (m *Mate) Name() string {
	return "Mate"
}

func (m *Mate) ExecutablesExists() bool {
	// Cache the results
	m.hasGsettings = which("gsettings") != ""
	m.hasMate = which("mate-session") != ""
	m.hasChecked = true

	// The result may be used both outside of this file, and in SetWallpaper
	return m.hasMate && m.hasGsettings
}

func (m *Mate) Running() bool {
	return containsE("GDMSESSION", "mate") || containsE("XDG_SESSION_DESKTOP", "MATE") || containsE("XDG_CURRENT_DESKTOP", "MATE") || containsE("DESKTOP_SESSION", "xsessions/mate")
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
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}
	// Check if dconf or gsettings are there, if we haven't already checked
	if !m.hasChecked {
		// This alters the state of m
		m.ExecutablesExists()
	}

	mode := defaultMode

	// If m.mode is specified, do not use the default value
	if m.mode != "" {
		mode = m.mode
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
		return fmt.Errorf("invalid desktop wallpaper picture mode for MATE: %s", mode)
	}

	if !m.hasGsettings {
		return errors.New("could not find gsettings")
	}

	//if MonConfOverlap("~/.config/monitors.xml") {
	//	return errors.New("there are overlapping monitor configurations in ~/.config/monitors.xml")
	//} else if m.verbose {
	//	fmt.Println("No monitor overlap in ~/.config/monitors.xml")
	//}

	// Exit if the monitor configuration will cause artifacts when setting
	// the desktop wallpaper.
	NoXrandrOverlapOrExit(m.verbose)

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
