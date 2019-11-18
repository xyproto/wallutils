// +build cgo

package xbg

// #cgo LDFLAGS: -lX11 -lXpm
// #include "xwallpaper.h"
import "C"
import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/xyproto/imagelib"
	"github.com/xyproto/xpm"
)

// X11 or Xorg windowmanager detector
type X11 struct {
	mode     string
	verbose  bool
	tempFile string
}

// Name returns the name of this window manager or desktop environment
func (x *X11) Name() string {
	return "X11"
}

// ExecutablesExists checks if executables associated with this backend exists in the PATH
func (x *X11) ExecutablesExists() bool {
	return which("X") != ""
}

// Running examines environment variables to try to figure out if either i3 or an X session is running (DISPLAY will then be set)
func (x *X11) Running() bool {
	// The X11 method of setting a wallpaper does not work correctly with some windowmanagers.

	// If NOFEH is set, the X11 backend will be used for i3 and AwesomeWM (useful for debugging the X11 backend)
	if hasE("NOFEH") {
		return hasE("DISPLAY")
	}

	// Check if i3 is running
	i3 := containsE("DESKTOP_SESSION", "i3") || containsE("XDG_CURRENT_DESKTOP", "i3") || containsE("XDG_SESSION_DESKTOP", "i3")

	// Check if AwesomeWM is running
	awm := containsE("DESKTOP_SESSION", "awesome") || containsE("XDG_CURRENT_DESKTOP", "awesome") || containsE("XDG_SESSION_DESKTOP", "awesome")

	// Recommend feh while working on the X11 backend
	if (i3 || awm) && (which("feh") == "") {
		// Output a helpful warning
		fmt.Println("Please install feh for i3 and AwesomeWM support")
	}

	// X is running, but not i3 or AwesomeWM
	return hasE("DISPLAY") && !(i3 || awm)
}

// SetMode will set the current way to display the wallpaper (stretched, tiled etc)
func (x *X11) SetMode(mode string) {
	x.mode = mode
}

// SetVerbose can be used for setting the verbose field to true or false.
// This will cause this backend to output information about what is is doing on stdout.
func (x *X11) SetVerbose(verbose bool) {
	x.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (x *X11) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}

	// Remove any existing temporary file before continuing
	if exists(x.tempFile) {
		if err := os.Remove(x.tempFile); err != nil {
			return err
		}
		x.tempFile = ""
	}

	// Generate a temporary filename
	tf, err := ioutil.TempFile("", "setwallpaper*.xpm")
	if err != nil {
		return err
	}
	convertedImageFilename := tf.Name()
	tf.Close()

	// Convert the given imageFilename to a temporary XPM file
	ext := strings.ToLower(filepath.Ext(imageFilename))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif":
		m, err := imagelib.Read(imageFilename)
		if err != nil {
			return err
		}
		imageName := filepath.Base(imageFilename[:len(imageFilename)-len(ext)])
		enc := xpm.NewEncoder(imageName)
		f, err := os.Create(convertedImageFilename)
		if err != nil {
			return err
		}
		defer f.Close()
		// Write the XPM image
		enc.Encode(f, m)
	default:
		return errors.New("unrecognized image file extension for: " + imageFilename)
	}

	if exists(convertedImageFilename) {
		imageFilename = convertedImageFilename
	} else {
		return errors.New("The generated XPM image does not exist: " + convertedImageFilename)
	}

	// Now that the file has been written, save the temporary filename for later deletion, at the next call to this function
	x.tempFile = imageFilename

	// NOTE: The C counterpart to this function may exit(1) if it's out of memory
	imageFilenameC := C.CString(imageFilename)

	// Should the C functions outut verbose messages to stdout?
	//verboseC := C.bool(x.verbose)
	verboseC := C.bool(true) // set to true while debugging an issue with the X11 backend

	// TODO: Figure out how to set the wallpaper mode
	retval := C.SetBackground(imageFilenameC, verboseC)

	C.free(unsafe.Pointer(imageFilenameC))
	switch retval {
	case -1:
		return errors.New("could not open X11 display with XOpenDisplay")
	}
	return nil
}
