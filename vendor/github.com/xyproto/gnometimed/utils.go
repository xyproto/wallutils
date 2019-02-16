package gnometimed

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// c formats a timestamp as HH:MM
func c(t time.Time) string {
	return fmt.Sprintf("%.2d:%.2d", t.Hour(), t.Minute())
}

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

// firstname finds the part of a filename before the extension
func firstname(filename string) string {
	ext := filepath.Ext(filename)
	return filename[:len(filename)-len(ext)]
}

// exists checks if the given path exists
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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

// CommonPrefix will find the longest common suffix in a slice of strings
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
// before the suffix. Will return the given string if it is too short to
// contain the prefix and suffix.
func Meat(s, prefix, suffix string) string {
	if len(s) < (len(prefix) + len(suffix)) {
		return s
	}
	return s[len(prefix) : len(s)-len(suffix)]
}
