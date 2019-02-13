// Package wallpaper provides a way to set the desktop wallpaper, for any windowmanager
package monitor

import (
	"errors"
	"fmt"
	"os"
)

type WM interface {
	Name() string
	ExecutablesExists() bool
	Running() bool
	SetWallpaper(string) error
	SetVerbose(bool)
}

type Wallpaper struct {
	CollectionName   string // the name of the directory containing this wallpaper, if it's not "pixmaps", "images" or "contents". May use the parent of the parent.
	Path             string // full path to the image filename
	Width            uint   // width of the image
	Height           uint   // height of the image
	PartOfCollection bool   // likely to be part of a wallpaper collection
}

// Default mode when setting the wallpaper for Gnome / Mate / Cinnamon
const defaultMode = "fill"

// All available backends for changing the wallpaper
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

func SetWallpaper(imageFilename string) error {
	return SetWallpaperVerbose(imageFilename, false)
}

// Set the desktop wallpaper (filled/stretched), for any supported windowmanager.
// The fallback is to use `feh`.
func SetWallpaperVerbose(imageFilename string, verbose bool) error {
	var lastErr error
	// Loop through all available WM structs
	for _, wm := range WMs {
		if wm.Running() && wm.ExecutablesExists() {
			if verbose {
				fmt.Printf("Using the %s backend.\n", wm.Name())
			}
			wm.SetVerbose(verbose)
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

func (wp *Wallpaper) String() string {
	star := " "
	if wp.PartOfCollection {
		star = "*"
	}
	return fmt.Sprintf("(%s) %dx%d\t%16s\t%s", star, wp.Width, wp.Height, wp.CollectionName, wp.Path)
}
