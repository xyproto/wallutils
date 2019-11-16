package wallutils

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// WM is an interface with the functions that needs to be implemented for adding support for setting the wallpaper for a new WM or DE
type WM interface {
	Name() string
	ExecutablesExists() bool
	Running() bool
	SetWallpaper(string) error
	SetVerbose(bool)
	SetMode(string)
}

// Wallpaper represents an image file that is part of a wallpaper collection (in a directory with several resolutions of the same image, for example)
type Wallpaper struct {
	CollectionName   string // the name of the directory containing this wallpaper, if it's not "pixmaps", "images" or "contents". May use the parent of the parent.
	Path             string // full path to the image filename
	Width            uint   // width of the image
	Height           uint   // height of the image
	PartOfCollection bool   // likely to be part of a wallpaper collection
}

// All backends should support these modes, if possible: stretch, fill, scale, tile, center
const defaultMode = "stretch"

// SetWallpaperCustom will set the given image filename as the wallpaper,
// regardless of which display server, window manager or desktop environment is in use.
func SetWallpaperCustom(imageFilename, mode string, verbose bool) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}
	var lastErr error
	// Loop through all available WM structs
	for _, wm := range WMs {
		if wm.Running() && wm.ExecutablesExists() {
			if verbose {
				fmt.Printf("Using the %s backend.\n", wm.Name())
			}
			wm.SetVerbose(verbose)
			if mode != "" && mode != defaultMode {
				wm.SetMode(mode)
			}
			if err := wm.SetWallpaper(imageFilename); err != nil {
				lastErr = err
				switch wm.Name() {
				case "Weston":
					// If the current windowmanager is Weston, no method is currently available
					return err
				default:
					if verbose {
						fmt.Fprintf(os.Stderr, "failed: %v\n", err)
					}
					// If the wallpaper mode is wrong, don't try the next backend, but return the error
					if strings.Contains(err.Error(), "invalid desktop wallpaper mode") {
						return err
					}
					// Try the next one
					continue
				}
			} else {
				return nil
			}
		}
	}
	if lastErr != nil {
		return fmt.Errorf("found no working method for setting the desktop wallpaper:\n%v", lastErr)
	}
	return errors.New("found no working method for setting the desktop wallpaper")

}

// SetWallpaperVerbose will set the desktop wallpaper, for any supported
// windowmanager. The fallback is to use `feh`. The wallpaper mode is "fill".
func SetWallpaperVerbose(imageFilename string, verbose bool) error {
	return SetWallpaperCustom(imageFilename, defaultMode, verbose)
}

// SetWallpaper will set the desktop wallpaper, for any supported
// windowmanager. The fallback is to use `feh`. The wallpaper mode is "fill".
func SetWallpaper(imageFilename string) error {
	return SetWallpaperCustom(imageFilename, defaultMode, false)
}

// Res returns the wallpaper resolution as a Res struct
func (wp *Wallpaper) Res() *Res {
	return NewRes(wp.Width, wp.Height)
}

// String returns a string with information about the wallpaper:
// - if it's part of a wallpaper collection or not
// - width
// - height
// - collection name
// - path
func (wp *Wallpaper) String() string {
	star := " "
	if wp.PartOfCollection {
		star = "*"
	}
	return fmt.Sprintf("(%s) %dx%d\t%16s\t%s", star, wp.Width, wp.Height, wp.CollectionName, wp.Path)
}
