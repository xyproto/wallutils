package monitor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// has checks if a string slice has the given element
func has(sl []string, e string) bool {
	for _, s := range sl {
		if s == e {
			return true
		}
	}
	return false
}

// unique removes all repeated elements from a slice of strings
func unique(sl []string) []string {
	var nl []string
	for _, s := range sl {
		if !has(nl, s) {
			nl = append(nl, s)
		}
	}
	return nl
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// firstname finds the part of a filename before the extension
func firstname(filename string) string {
	ext := filepath.Ext(filename)
	return filename[:len(filename)-len(ext)]
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func which(executable string) string {
	p, err := exec.LookPath(executable)
	if err != nil {
		return ""
	}
	return p
}

func run(shellCommand string) error {
	fmt.Println(shellCommand)
	cmd := exec.Command("sh", "-c", shellCommand)
	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}
	return nil
}

func output(shellCommand string) string {
	fmt.Println(shellCommand)
	cmd := exec.Command("sh", "-c", shellCommand)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return string(stdoutStderr)
}

func containsE(envVar, subString string) bool {
	return strings.Contains(os.Getenv(envVar), subString)
}

func hasE(envVar string) bool {
	return os.Getenv(envVar) != ""
}
