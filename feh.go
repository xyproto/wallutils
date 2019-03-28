package wallutils

// This is the fallback if no specific windowmanager has been detected

import (
	"errors"
	"fmt"
	"strings"
)

type Feh struct {
	mode    string
	verbose bool
}

func (f *Feh) Name() string {
	return "Feh"
}

func (f *Feh) ExecutablesExists() bool {
	return which("feh") != ""
}

func (f *Feh) Running() bool {
	return true
}

func (f *Feh) SetMode(mode string) {
	f.mode = mode
}

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
	// if missing, prefix with "bg-"
	if !strings.HasPrefix(mode, "bg-") {
		mode = "bg-" + mode
	}

	// check if the mode is valid
	switch mode {
	case "bg-fill", "bg-center", "bg-max", "bg-scale", "bg-tile":
		break
	default:
		// Invalid and unrecognized desktop wallpaper mode
		return fmt.Errorf("invalid desktop wallpaper mode for Feh: %s", mode)
	}

	// set the wallpaper with feh
	if err := run("feh", []string{"--" + mode, imageFilename}, f.verbose); err != nil {
		return errors.New("feh --" + mode + " " + imageFilename + " failed to run")
	}
	return nil
}
