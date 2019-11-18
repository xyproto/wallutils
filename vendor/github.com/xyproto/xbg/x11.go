// +build cgo

package xbg

// #cgo CFLAGS: -DXINERAMA
// #cgo LDFLAGS: -lX11 -lImlib2 -lXinerama -lm
// #include "bgs.h"
import "C"
import (
	"errors"
	"fmt"
	"sync"
	"unsafe"
)

var singletonComplete bool
var singletonMut = &sync.RWMutex{}

// X11 or Xorg windowmanager detector
type X11 struct {
	mode    string
	verbose bool
	rotate  bool
	mut     *sync.RWMutex
}

// New will only create an X11 struct once per run.
// It's a singleton.
func New() *X11 {
	singletonMut.RLock()
	if singletonComplete {
		singletonMut.RUnlock()
		return nil
	}
	singletonMut.RUnlock()
	singletonMut.Lock()
	singletonComplete = true
	singletonMut.Unlock()
	m := &sync.RWMutex{}
	return &X11{"scale", false, false, m}
}

// Release will set this struct to nil and make it possible to aquire a new one later on
func (x *X11) Release() {
	singletonMut.RLock()
	if x != nil {
		x = nil
	} else {
		singletonMut.RUnlock()
		return
	}
	singletonMut.RUnlock()
	singletonMut.RLock()
	if singletonComplete {
		singletonMut.RUnlock()
		singletonMut.Lock()
		singletonComplete = false
		singletonMut.Unlock()
		return
	}
	singletonMut.RUnlock()
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
	// Check if i3 is running
	//i3 := containsE("DESKTOP_SESSION", "i3") || containsE("XDG_CURRENT_DESKTOP", "i3") || containsE("XDG_SESSION_DESKTOP", "i3")

	// Check if AwesomeWM is running
	//awm := containsE("DESKTOP_SESSION", "awesome") || containsE("XDG_CURRENT_DESKTOP", "awesome") || containsE("XDG_SESSION_DESKTOP", "awesome")

	// Just check if DISPLAY is set
	return hasE("DISPLAY")
}

// SetMode will set the current way to display the wallpaper (stretched, tiled etc)
func (x *X11) SetMode(mode string) {
	x.mut.Lock()
	x.mode = mode
	x.mut.Unlock()
}

// SetVerbose can be used for setting the verbose field to true or false.
// This will cause this backend to output information about what is is doing on stdout.
func (x *X11) SetVerbose(verbose bool) {
	x.mut.Lock()
	x.verbose = verbose
	x.mut.Unlock()
}

// SetRotate can be used for setting the rotate field to true or false.
func (x *X11) SetRotate(rotate bool) {
	x.mut.Lock()
	x.rotate = rotate
	x.mut.Unlock()
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (x *X11) SetWallpaper(imageFilename string) error {
	x.mut.Lock()

	if !exists(imageFilename) {
		x.mut.Unlock()
		return fmt.Errorf("no such file: %s", imageFilename)
	}

	var mode C.ImageMode
	switch x.mode {
	case "center":
		mode = C.ImageMode(C.ModeCenter)
	case "zoom", "zoomed", "fill", "max", "":
		mode = C.ImageMode(C.ModeZoom)
	case "scale", "scaled", "stretch", "stretched":
		mode = C.ImageMode(C.ModeScale)
	default:
		x.mut.Unlock()
		// for unsupported modes: "fit", "tiled" or anything
		return fmt.Errorf("unsupported desktop wallpaper mode for the X11 backend: %s", x.mode)
	}

	x.mut.Unlock()
	x.mut.RLock()

	// Prepare a the image filename as a char* string
	cFilename := C.CString(imageFilename)

	// Set the background image, and intepret the returned string as a Go string
	errString := C.GoString(C.SetBackground(cFilename, C.bool(x.rotate), mode, C.bool(x.verbose)))

	// Free the C string
	C.free(unsafe.Pointer(cFilename))

	x.mut.RUnlock()

	// Return an error, if errString is not empty
	if len(errString) > 0 {
		return errors.New("could not open X11 display with XOpenDisplay: " + errString)
	}
	return nil
}
