package wallutils

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os/user"
	"path"
	"strings"
	"time"

	"github.com/xyproto/env/v2"
)

// Hyprpaper compatible windowmanager
type Hyprpaper struct {
	mode    string
	sock    string
	verbose bool
}

// Name returns the name of this window manager or desktop environment
func (hp *Hyprpaper) Name() string {
	return "Hyprpaper"
}

// ExecutablesExists checks if executables associated with this backend exists in the PATH
func (hp *Hyprpaper) ExecutablesExists() bool {
	return which("hyprpaper") != "" // && which("hyprctl") != ""
}

// Running examines environment variables to try to figure out if this backend is currently running.
// Also sets hp.sock to an empty string or to a UNIX socket file.
func (hp *Hyprpaper) Running() bool {
	hp.sock = ""

	inst := env.Str("HYPRLAND_INSTANCE_SIGNATURE")
	if inst == "" {
		return false // HYPRLAND_INSTANCE_SIGNATURE is not set
	}

	currentUser, err := user.Current()
	if err != nil || currentUser.Uid == "" {
		return false // could not get the user ID of the current user
	}

	if sock := path.Join("/run/user/" + currentUser.Uid + "/hypr/" + inst + "/.hyprpaper.sock"); exists(sock) {
		hp.sock = sock
		return true
	}

	if sock := path.Join("/tmp/hypr", inst, ".hyprpaper.sock"); exists(sock) {
		hp.sock = sock
		return true
	}

	return false
}

// SetMode will set the current way to display the wallpaper (stretched, tiled etc)
func (hp *Hyprpaper) SetMode(mode string) {
	hp.mode = mode
}

// SetVerbose can be used for setting the verbose field to true or false.
// This will cause this backend to output information about what is is doing on stdout.
func (hp *Hyprpaper) SetVerbose(verbose bool) {
	hp.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (hp *Hyprpaper) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}

	// Connect to the Hyprpaper socket
	sock, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: hp.sock, Net: "unix"})

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
	if hp.mode != "" {
		mode = hp.mode
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
	if _, err := sock.Write([]byte(cmd)); err != nil {
		return err
	}

	sock.SetReadDeadline(time.Now().Add(2 * time.Second))

	data, err := io.ReadAll(sock)
	if err != nil {
		var ne net.Error
		if !(errors.As(err, &ne) && ne.Timeout()) {
			return err
		}
	}

	res := strings.TrimSpace(string(data))
	if res != "ok" {
		return fmt.Errorf("cmd %q failed: %s", cmd, res)
	}

	return nil
}
