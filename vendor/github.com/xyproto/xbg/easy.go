package xbg

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// SetWallpaper will set the given image filename as the wallpaper,
// The mode string can be:
//   "center" for "center mode"
//   "zoom", "zoomed", "fill" or "max" for "zoom mode"
//   "scale", "scaled", "stretch" or "streatched" for "scale mode"
// If verbose is set, some output will be written to stdout.
func SetWallpaper(imageFilename, mode string, verbose bool) error {
	// This check is redundant, but it's nice to check it before checking if the WM is ready
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}
	wm := New()
	defer wm.Release()
	if wm == nil {
		return errors.New("the X11 wallpaper backend is already in use")
	}
	if !(wm.Running() && wm.ExecutablesExists()) {
		return errors.New("found no working method for setting the desktop wallpaper, maybe X11 is not installed or DISPLAY not set")
	}
	if verbose {
		fmt.Printf("Using the %s backend.\n", wm.Name())
	}
	wm.verbose = verbose
	wm.mode = mode
	if err := wm.SetWallpaper(imageFilename); err != nil {
		if verbose {
			fmt.Fprintf(os.Stderr, "failed: %v\n", err)
		}
		// If the wallpaper mode is wrong, don't try the next backend, but return the error
		if strings.Contains(err.Error(), "invalid desktop wallpaper mode") {
			return err
		}
		// This did not work out
		return errors.New("could not use the X11 backend")
	}
	return nil
}
