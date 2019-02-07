package monitor

// This is the fallback if no specific windowmanager has been detected

import (
	"errors"
)

type Feh struct {
}

func (g *Feh) Name() string {
	return "Feh"
}

func (g *Feh) ExecutablesExists() bool {
	return which("feh") != ""
}

func (g *Feh) Running() bool {
	return true
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
// `feh` is used for setting the desktop background, and must be in the PATH.
func (g *Feh) SetWallpaper(imageFilename string) error {
	// bg-fill | bg-center | bg-max | bg-scale | bg-tile
	command := "feh --bg-fill " + imageFilename
	if err := run(command); err != nil {
		return errors.New(command + " failed to run")
	}
	return nil
}
