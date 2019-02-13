package monitor

// #cgo LDFLAGS: -lX11
// #include "xwallpaper.h"
import "C"
import (
	"errors"
	"unsafe"
)

// X11 or Xorg windowmanager detector
type X11 struct {
	verbose bool
}

func (x *X11) Name() string {
	return "X11"
}

func (x *X11) ExecutablesExists() bool {
	return which("X") != ""
}

func (x *X11) Running() bool {
	// The X11 method of setting a wallpaper does not seem to work with i3
	i3 := containsE("DESKTOP_SESSION", "i3") || containsE("XDG_CURRENT_DESKTOP", "i3") || containsE("XDG_SESSION_DESKTOP", "i3")
	return !i3 && hasE("DISPLAY")
}

func (x *X11) SetVerbose(verbose bool) {
	x.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (s *X11) SetWallpaper(imageFilename string) error {
	// NOTE: The C counterpart to this function may exit(1) if it's out of memory
	imageFilenameC := C.CString(imageFilename)
	retval := C.SetBackground(imageFilenameC)
	C.free(unsafe.Pointer(imageFilenameC))
	switch retval {
	case -1:
		return errors.New("could not open X11 display with XOpenDisplay")
	}
	return nil
}
