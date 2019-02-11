// Package wallpaper provides a way to set the desktop wallpaper, for any windowmanager
package monitor

import (
	"errors"
	"fmt"
)

type WM interface {
	Name() string
	ExecutablesExists() bool
	Running() bool
	SetWallpaper(imageFilename string) error
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
	&X11{}, // using a C+Go .so plugin
}

var (
	DefaultWallpaperDirectories = []string{"/usr/share/pixmaps", "/usr/share/wallpapers", "/usr/share/backgrounds", "/usr/local/share/pixmaps", "/usr/local/share/wallpapers", "/usr/local/share/backgrounds", "/usr/share/archlinux"}
)

// Set the desktop wallpaper (filled/stretched), for any supported windowmanager.
// The fallback is to use `feh`.
func SetWallpaper(imageFilename string) error {
	var lastErr error
	// Loop through all available WM structs
	for _, wm := range WMs {
		if wm.Running() && wm.ExecutablesExists() {
			fmt.Printf("Setting wallpaper with the %s backend.\n", wm.Name())
			if err := wm.SetWallpaper(imageFilename); err != nil {
				lastErr = err
				switch wm.Name() {
				case "Weston":
					// If the current windowmanager is Weston, no method is currently available
					return err
				default:
					fmt.Println("failed:", err)
					// Try the next one
					continue
				}
			} else {
				return nil
			}
		}
	}
	if lastErr != nil {
		return fmt.Errorf("Found no working method for setting the desktop wallpaper:\n%s", lastErr)
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
