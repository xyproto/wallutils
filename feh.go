package wallutils

// This is the fallback if no specific windowmanager has been detected

import (
	"errors"
	"fmt"
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
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}
	// bg-fill | bg-center | bg-max | bg-scale | bg-tile
	mode := "bg-fill"
	if err := run("feh", []string{"--" + mode, imageFilename}, f.verbose); err != nil {
		return errors.New("feh --" + mode + " " + imageFilename + " failed to run")
	}
	return nil
}
