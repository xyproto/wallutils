package wallutils

import (
	"fmt"
	"net"
	"path"

	"github.com/xyproto/env/v2"
)

// Hyprpaper compatible windowmanager
type Hyprpaper struct {
	mode    string
	sock    string
	verbose bool
}

// Name returns the name of this window manager or desktop environment
func (sb *Hyprpaper) Name() string {
	return "Hyprpaper"
}

// ExecutablesExists checks if executables associated with this backend exists in the PATH
func (sb *Hyprpaper) ExecutablesExists() bool {
	return which("hyprpaper") != ""
}

// Running examines environment variables to try to figure out if this backend is currently running
func (sb *Hyprpaper) Running() bool {
	inst := env.Str("HYPRLAND_INSTANCE_SIGNATURE")
	sb.sock = path.Join("/tmp/hypr", inst, ".hyprpaper.sock")
	return exists(sb.sock)
}

// SetMode will set the current way to display the wallpaper (stretched, tiled etc)
func (sb *Hyprpaper) SetMode(mode string) {
	sb.mode = mode
}

// SetVerbose can be used for setting the verbose field to true or false.
// This will cause this backend to output information about what is is doing on stdout.
func (sb *Hyprpaper) SetVerbose(verbose bool) {
	sb.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (sb *Hyprpaper) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}

	// Connect to the Hyprpaper socket
	sock, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: sb.sock, Net: "unix"})

	if err != nil {
		return err
	}

	defer sock.Close()

	// Preload the image
	if err = runHyprCmd(sock, "preload "+imageFilename); err != nil {
		return err
	}

	// Set the wallpaper mode
	mode := defaultMode
	if sb.mode != "" {
		mode = sb.mode
	}

	switch mode {
	case "center", "tile", "fill", "stretch", "scale", "scaled", "zoom", "zoomed", "stretched":
		mode = ""
	case "fit":
		mode = "contain:"
	default:
		// Invalid and unrecognized desktop wallpaper mode
		return fmt.Errorf("invalid desktop wallpaper mode for swaybg: %s", mode)
	}

	// Set the wallpaper
	if err = runHyprCmd(sock, "wallpaper ,"+mode+imageFilename); err != nil {
		return err
	}

	// Unload unused images
	return runHyprCmd(sock, "unload unused")
}

func runHyprCmd(sock *net.UnixConn, cmd string) error {
	_, err := sock.Write([]byte(cmd))

	if err != nil {
		return err
	}

	buf := [32]byte{}

	n, err := sock.Read(buf[:])

	if err != nil {
		return err
	}

	res := string(buf[:n])

	if res != "ok" {
		return fmt.Errorf("cmd %q failed: %s", cmd, res)
	}

	return nil
}
