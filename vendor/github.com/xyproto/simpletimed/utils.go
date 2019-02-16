package simpletimed

import (
	"fmt"
	"os"
	"time"
)

var h24 = time.Hour * 24

// c formats a timestamp as HH:MM
func c(t time.Time) string {
	return fmt.Sprintf("%.2d:%.2d", t.Hour(), t.Minute())
}

// exists checks if the given path exists
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
