package monitor

// This is the fallback if no specific windowmanager has been detected

import (
	"errors"
)

type Feh struct {
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

func (f *Feh) SetVerbose(verbose bool) {
	f.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
// `feh` is used for setting the desktop background, and must be in the PATH.
func (f *Feh) SetWallpaper(imageFilename string) error {
	// bg-fill | bg-center | bg-max | bg-scale | bg-tile
	command := "feh --bg-fill " + imageFilename
	if err := run(command, f.verbose); err != nil {
		return errors.New(command + " failed to run")
	}
	return nil
}
