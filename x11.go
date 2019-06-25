// +build cgo

package wallutils

// #cgo LDFLAGS: -lX11 -lXpm
// #include "xwallpaper.h"
import "C"
import (
	"errors"
	"fmt"
	"github.com/xyproto/imagelib"
	"github.com/xyproto/xpm"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

// X11 or Xorg windowmanager detector
type X11 struct {
	mode    string
	verbose bool
}

func (x *X11) Name() string {
	return "X11"
}

func (x *X11) ExecutablesExists() bool {
	return which("X") != ""
}

func (x *X11) Running() bool {
	// The X11 method of setting a wallpaper does not seem to work with i3,
	// so check if i3 is running first.
	i3 := containsE("DESKTOP_SESSION", "i3") || containsE("XDG_CURRENT_DESKTOP", "i3") || containsE("XDG_SESSION_DESKTOP", "i3")

	// X is running, but not i3
	return hasE("DISPLAY") && !i3
}

func (x *X11) SetMode(mode string) {
	x.mode = mode
}

func (x *X11) SetVerbose(verbose bool) {
	x.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (x *X11) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}

	convertedImageFilename := filepath.Join("/tmp", "_setwallpaper.xpm")
	ext := strings.ToLower(filepath.Ext(imageFilename))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif":
		m, err := imagelib.Read(imageFilename)
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
	}
	// TODO: Return an error here if convertedImageFilename does not exist,
	//       or return an error later, when trying to set the imageFilename?

	// NOTE: The C counterpart to this function may exit(1) if it's out of memory
	imageFilenameC := C.CString(imageFilename)
	// TODO: Figure out how to set the wallpaper mode
	retval := C.SetBackground(imageFilenameC)
	C.free(unsafe.Pointer(imageFilenameC))
	switch retval {
	case -1:
		return errors.New("could not open X11 display with XOpenDisplay")
	}
	return nil
}
