package wallutils

import (
	"strings"
)

// sq will single-quote a string
func sq(s string) string {
	return "'" + s + "'"
}

// GSettings can be used for getting and setting configuration options with gsettings
type GSettings struct {
	schema  string
	verbose bool
}

// NewGSettings creates a new GSettings struct given a schema/category and a
// bool for if the commands should be printed to stdout before running.
func NewGSettings(schema string, verbose bool) *GSettings {
	return &GSettings{schema: schema, verbose: verbose}
}

// Set a key using gsettings. Will single-quote the given value.
func (g *GSettings) Set(key, value string) error {
	return run("gsettings", []string{"set", g.schema, key, sq(value)}, g.verbose)
}

// Get a key using gsettings. Will return an empty string if there are errors.
func (g *GSettings) Get(key string) string {
	retval := strings.TrimSpace(output("gsettings", []string{"get", g.schema, key}, g.verbose))
	// Remove single quotes, if present
	if strings.HasPrefix(retval, "'") && strings.HasSuffix(retval, "'") {
		return retval[1 : len(retval)-1]
	}
	return retval
}
