package wallutils

// This is for Pekwm 0.2.0 or later

import (
	"errors"
	"fmt"

	"github.com/xyproto/env/v2"
)

// Pekwm is a structure containing settings for running the "pekwm_bg" executable
type Pekwm struct {
	mode    string
	verbose bool
}

// Name returns the name of this method of setting a wallpaper
func (f *Pekwm) Name() string {
	return "Pekwm"
}

// ExecutablesExists checks if the "pekwm_bg" executable exists in the PATH
// (comes with pekwm 0.2.0 or later)
func (f *Pekwm) ExecutablesExists() bool {
	return which("pekwm_bg") != ""
}

// Running checks if $PEKWM_CONFIG_FILE is set
func (f *Pekwm) Running() bool {
	return env.Has("PEKWM_CONFIG_FILE")
}

// SetMode will set the current way to display the wallpaper (stretched, tiled etc).
// The selected mode must be compatible with pekwm_bg.
func (f *Pekwm) SetMode(mode string) {
	f.mode = mode
}

// SetVerbose can be used for setting the verbose field to true or false.
// This will cause this backend to output information about what is is doing on stdout.
func (f *Pekwm) SetVerbose(verbose bool) {
	f.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (f *Pekwm) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}
	mode := defaultMode
	if f.mode != "" {
		mode = f.mode
	}

	// Images are set with ie. "pekwm_bg Image /somewhere/image.png#scaled"

	var tag string

	// check if the mode is valid
	switch mode {
	case "stretch", "stretched", "scaled", "fill", "scale", "max":
		tag = "#scaled"
	case "tile", "tiled":
		// pekwm_bg tiles by default
		break
	default:
		// Invalid and unrecognized desktop wallpaper mode
		return fmt.Errorf("invalid desktop wallpaper mode for Pekwm: %s", mode)
	}

	// set the wallpaper with pekwm_bg
	if err := run("pekwm_bg", []string{"-D", "Image", imageFilename + tag}, f.verbose); err != nil {
		return errors.New("pekwm_bg -D Image \"" + imageFilename + tag + "\" failed to run")
	}
	return nil
}
