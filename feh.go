package wallutils

// This is the fallback if no specific windowmanager has been detected

import (
	"errors"
	"fmt"
	"strings"
)

// Feh is a structure containing settings for running the "feh" executble
type Feh struct {
	mode    string
	verbose bool
}

// Name returns the name of this method of setting a wallpaper
func (f *Feh) Name() string {
	return "Feh"
}

// ExecutablesExists checks if the feh executable exists in the PATH
func (f *Feh) ExecutablesExists() bool {
	return which("feh") != ""
}

// Running just returns true for the Feh backend, since this is an application and not a WM / DM
func (f *Feh) Running() bool {
	return true
}

// SetMode will set the current way to display the wallpaper (stretched, tiled etc).
// The selected mode must be compatible with feh.
func (f *Feh) SetMode(mode string) {
	f.mode = mode
}

// SetVerbose can be used for setting the verbose field to true or false.
// This will cause this backend to output information about what is is doing on stdout.
func (f *Feh) SetVerbose(verbose bool) {
	f.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
// `feh` is used for setting the desktop background, and must be in the PATH.
func (f *Feh) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}
	mode := defaultMode
	if f.mode != "" {
		mode = f.mode
	}
	// remove the "bg-" prefix, if it's there
	if strings.HasPrefix(mode, "bg-") {
		mode = mode[3:]
	}

	// check if the mode is valid
	switch mode {
	case "fill", "center", "max", "scale", "tile":
		break
	case "zoom", "zoomed":
		mode = "fill"
	case "stretch", "stretched", "scaled":
		mode = "scale"
	case "fit":
		mode = "max"
	default:
		// Invalid and unrecognized desktop wallpaper mode
		return fmt.Errorf("invalid desktop wallpaper mode for Feh: %s", mode)
	}

	// set the wallpaper with feh
	if err := run("feh", []string{"--bg-" + mode, imageFilename}, f.verbose); err != nil {
		return errors.New("feh --bg-" + mode + " " + imageFilename + " failed to run")
	}
	return nil
}
