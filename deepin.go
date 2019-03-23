package wallutils

import (
	"errors"
	"fmt"
)

// Deepin windowmanager detector
type Deepin struct {
	mode         string // none | wallpaper | centered | scaled | stretched | zoom | spanned, scaled is the default
	hasDeepin    bool
	hasGsettings bool
	hasChecked   bool
	verbose      bool
}

func (d *Deepin) Name() string {
	return "Deepin"
}

func (d *Deepin) ExecutablesExists() bool {
	// Cache the results
	d.hasDeepin = which("deepin-session") != ""
	d.hasGsettings = which("gsettings") != ""
	d.hasChecked = true

	// The result may be used both outside of this file, and in SetWallpaper
	return d.hasDeepin && d.hasGsettings
}

func (d *Deepin) Running() bool {
	return containsE("GDMSESSION", "deepin") || containsE("XDG_CURRENT_DESKTOP", "Deepin") || containsE("DESKTOP_SESSION", "xsessions/deepin")
}

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

	mode := defaultMode

	// If d.mode is specified, do not use the default value
	if d.mode != "" {
		mode = d.mode
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
		return fmt.Errorf("invalid desktop wallpaper picture mode for Deepin: %s", mode)
	}

	if !d.hasGsettings {
		return errors.New("could not find gsettings")
	}

	//if MonConfOverlap("~/.config/monitors.xml") {
	//	return errors.New("there are overlapping monitor configurations in ~/.config/monitors.xml")
	//} else if d.verbose {
	//	fmt.Println("No monitor overlap in ~/.config/monitors.xml")
	//}

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

	// Set the desktop wallpaper (also set it if it is already set, in case
	// the contents have changed)
	return g.Set("picture-uri", "file://"+imageFilename)
}
