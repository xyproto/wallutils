package xbg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

func containsE(envVar, subString string) bool {
	return strings.Contains(os.Getenv(envVar), subString)
}

func hasE(envVar string) bool {
	return os.Getenv(envVar) != ""
}

// CommonPrefix will find the longest common prefix in a slice of strings
func CommonPrefix(sl []string) string {
	if len(sl) == 0 {
		return ""
	}
	shortestLength := len(sl[0])
	shortestString := sl[0]
	for _, s := range sl {
		if len(s) < shortestLength {
			shortestLength = len(s)
			shortestString = s
		}
	}
	if shortestLength == 0 {
		return ""
	}
	for i := 1; i < shortestLength; i++ {
		for _, s := range sl {
			if !strings.HasPrefix(s, shortestString[:i]) {
				return shortestString[:i-1]
			}
		}
	}
	return shortestString
}

// CommonSuffix will find the longest common suffix in a slice of strings
func CommonSuffix(sl []string) string {
	if len(sl) == 0 {
		return ""
	}
	shortestLength := len(sl[0])
	shortestString := sl[0]
	for _, s := range sl {
		if len(s) < shortestLength {
			shortestLength = len(s)
			shortestString = s
		}
	}
	if shortestLength == 0 {
		return ""
	}
	for i := 1; i < shortestLength; i++ {
		for _, s := range sl {
			if !strings.HasSuffix(s, shortestString[shortestLength-i:]) {
				return shortestString[shortestLength-(i-1):]
			}
		}
	}
	return shortestString
}

// Meat returns the meat of the string: the part that is after the prefix and
// before the suffix. It does not check if the prefix and suffix exists in the
// string. If the given string is too short to contain the prefix and suffix,
// it will be returned as it is.
func Meat(s, prefix, suffix string) string {
	if len(s) < (len(prefix) + len(suffix)) {
		return s
	}
	return s[len(prefix) : len(s)-len(suffix)]
}

// Quit with a nicely formatted error message to stderr
func Quit(err error) {
	msg := err.Error()
	if !strings.HasSuffix(msg, ".") && !strings.HasSuffix(msg, "!") && !strings.Contains(msg, ":") {
		msg += "."
	}
	fmt.Fprintf(os.Stderr, "%s%s\n", strings.ToUpper(string(msg[0])), msg[1:])
	os.Stdout.Sync()
	os.Stderr.Sync()
	os.Exit(1)
}
