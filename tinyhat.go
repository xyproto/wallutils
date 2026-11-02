package wallutils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"

	"github.com/xyproto/env/v2"
)

// dwl/Tiny Hat
type TinyHat struct {
	mode    string
	verbose bool
}

// Name returns the name of this window manager or desktop environment
func (th *TinyHat) Name() string {
	return "Tiny Hat"
}

// ExecutablesExists checks if executables associated with this backend exists in the PATH
func (th *TinyHat) ExecutablesExists() bool {
	return which("dwl") != ""
}

// Running checks if the Tiny Hat dwl session is running and reading commands
func (th *TinyHat) Running() bool {
	tinyHatPipe := filepath.Join(env.Str("XDG_RUNTIME_DIR", "/tmp"), "tinyhat-bar")
	fi, err := os.Stat(tinyHatPipe)
	return err == nil && fi.Mode()&os.ModeNamedPipe != 0 && env.Contains("XDG_CURRENT_DESKTOP", "dwl")
}

// SetMode will set the current way to display the wallpaper (stretched, tiled etc)
func (th *TinyHat) SetMode(mode string) {
	th.mode = mode
}

// SetVerbose can be used for setting the verbose field to true or false.
// This will cause this backend to output information about what is is doing on stdout.
func (th *TinyHat) SetVerbose(verbose bool) {
	th.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given a PNG image filename.
// The image must exist and be readable.
func (th *TinyHat) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}
	path, err := filepath.Abs(imageFilename)
	if err != nil {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	header := make([]byte, 8)
	_, err = io.ReadFull(f, header)
	f.Close()
	if err != nil || !bytes.Equal(header, []byte("\x89PNG\r\n\x1a\n")) {
		return fmt.Errorf("dwl/Tiny Hat can only show PNG images: %s", imageFilename)
	}

	mode := defaultMode
	if th.mode != "" {
		mode = th.mode
	}

	switch mode {
	case "center", "tile", "fill", "stretch", "fit", "left":
	case "scale", "scaled":
		mode = "fill"
	case "zoom", "zoomed", "stretched":
		mode = "stretch"
	default:
		return fmt.Errorf("invalid desktop wallpaper mode for dwl/Tiny Hat: %s", mode)
	}

	tinyHatPipe := filepath.Join(env.Str("XDG_RUNTIME_DIR", "/tmp"), "tinyhat-bar")
	pipe, err := os.OpenFile(tinyHatPipe, os.O_WRONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return fmt.Errorf("could not reach dwl: %w", err)
	}
	defer pipe.Close()

	command := fmt.Sprintf("helptext none\nwallpaper %s %s\n", mode, path)
	if th.verbose {
		fmt.Print(command)
	}
	_, err = pipe.Write([]byte(command))
	return err
}
