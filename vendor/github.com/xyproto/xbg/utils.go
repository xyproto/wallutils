package xbg

import (
	"os"
	"os/exec"
)

// exists checks if the given path exists
func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// which tries to find the given executable name in the $PATH
// Returns an empty string if not found.
func which(executable string) string {
	p, err := exec.LookPath(executable)
	if err != nil {
		return ""
	}
	return p
}

func hasE(envVar string) bool {
	return os.Getenv(envVar) != ""
}
