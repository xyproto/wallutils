package simpletimed

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var h24 = time.Hour * 24

// cFmt formats a timestamp as HH:MM
func cFmt(t time.Time) string {
	return fmt.Sprintf("%.2d:%.2d", t.Hour(), t.Minute())
}

// dFmt formats a duration nicely
func dFmt(d time.Duration) string {
	s := fmt.Sprintf("%6s", d)
	if strings.Contains(s, ".") {
		pos := strings.Index(s, ".")
		if strings.Contains(s[pos:], "s") {
			s = s[:pos] + "s"
		}
	}
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return strings.TrimSpace(s)
}

// exists checks if the given path exists
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// mod24 returns the duration, but in the interval from 0 to 24, regardless
// if the given duration is negative, less than 24h or larger than 24h.
// The modulo operator is used to "wrap" the given duration in a 24h interval.
// Unlike in Python, modulo in Go can return a negative result.
func mod24(d time.Duration) time.Duration {
	hourDiff := d % h24
	if hourDiff < 0 {
		return -hourDiff
	}
	return hourDiff
}
