//go:build cgo
// +build cgo

package wallutils

import (
	"fmt"

	"github.com/xyproto/env/v2"
)

// SwayBG compatible windowmanager
type SwayBG struct {
	mode    string
	verbose bool
}

// Name returns the name of this window manager or desktop environment
func (sb *SwayBG) Name() string {
	return "SwayBG"
}

// ExecutablesExists checks if executables associated with this backend exists in the PATH
func (sb *SwayBG) ExecutablesExists() bool {
	return which("swaybg") != ""
}

// Running examines environment variables to try to figure out if this backend is currently running
func (sb *SwayBG) Running() bool {
	return env.Has("WAYLAND_DISPLAY") || env.Str("XDG_SESSION_TYPE") == "wayland"
}

// SetMode will set the current way to display the wallpaper (stretched, tiled etc)
func (sb *SwayBG) SetMode(mode string) {
	sb.mode = mode
}

// SetVerbose can be used for setting the verbose field to true or false.
// This will cause this backend to output information about what is is doing on stdout.
func (sb *SwayBG) SetVerbose(verbose bool) {
	sb.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (sb *SwayBG) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}

	// initialize the mode setting (stretched/tiled etc)
	mode := defaultMode
	if sb.mode != "" {
		mode = sb.mode
	}

	// stretch, fit, fill, center, tile, or solid_color
	switch mode {
	case "center", "tile", "fill", "stretch", "fit":
		break
	case "scale", "scaled":
		mode = "fill"
	case "zoom", "zoomed", "stretched":
		mode = "stretch"
	default:
		// Invalid and unrecognized desktop wallpaper mode
		return fmt.Errorf("invalid desktop wallpaper mode for swaybg: %s", mode)
	}

	// first stop swaybg, if it`s already running (and ignore errors, if it can't be killed)
	run("pkill", []string{"swaybg"}, sb.verbose)

	// start a new instance
	pid, err := runbg("swaybg", []string{"-i", imageFilename, "-m", mode, "-o", "*"}, sb.verbose)
	if err != nil {
		return err
	}

	if sb.verbose {
		// output the new PID
		fmt.Println("started PID", pid)
	}
	return nil
}
