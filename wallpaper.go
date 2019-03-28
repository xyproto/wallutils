package wallutils

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type WM interface {
	Name() string
	ExecutablesExists() bool
	Running() bool
	SetWallpaper(string) error
	SetVerbose(bool)
	SetMode(string)
}

type Wallpaper struct {
	CollectionName   string // the name of the directory containing this wallpaper, if it's not "pixmaps", "images" or "contents". May use the parent of the parent.
	Path             string // full path to the image filename
	Width            uint   // width of the image
	Height           uint   // height of the image
	PartOfCollection bool   // likely to be part of a wallpaper collection
}

// "fill" is not supported by all DEs and WMs, but all backends here should
// either support "fill" or use something equivalent.
const defaultMode = "fill"

// WMs contains all available backends for changing the wallpaper
var WMs = []WM{
	&Sway{},
	&Deepin{},
	&Xfce4{},
	&Mate{},
	&Cinnamon{},
	&Plasma{},
	&Gnome3{},
	&Gnome2{},
	&Weston{},
	&Feh{}, // using feh
	&X11{},
}

func SetWallpaperCustom(imageFilename string, verbose bool, mode string) error {
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
		return fmt.Errorf("Found no working method for setting the desktop wallpaper:\n%v", lastErr)
	}
	// Flush stderr
	if err := os.Stderr.Sync(); err != nil {
		return err
	}
	// Flush stdout
	if err := os.Stdout.Sync(); err != nil {
		return err
	}
	return errors.New("Found no working method for setting the desktop wallpaper")

}

// SetWallpaperVerbose will set the desktop wallpaper, for any supported
// windowmanager. The fallback is to use `feh`. The wallpaper mode is "fill".
func SetWallpaperVerbose(imageFilename string, verbose bool) error {
	return SetWallpaperCustom(imageFilename, verbose, defaultMode)
}

// SetWallpaper will set the desktop wallpaper, for any supported
// windowmanager. The fallback is to use `feh`. The wallpaper mode is "fill".
func SetWallpaper(imageFilename string) error {
	return SetWallpaperCustom(imageFilename, false, defaultMode)
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
