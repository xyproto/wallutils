package monitor

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

// Res is a structure containing width and height
type Res struct {
	w, h uint
}

// NewRes creates a new resolution structure
func NewRes(w, h uint) *Res {
	return &Res{w, h}
}

func (r *Res) String() string {
	return fmt.Sprintf("%dx%d", r.w, r.h)
}

// W is the width
func (r *Res) W() uint {
	return r.w
}

// H is the height
func (r *Res) H() uint {
	return r.h
}

// Distance returns the distance between two resolutions (Euclidean distance)
func Distance(a, b *Res) int {
	return abs(int(b.w)-int(a.w)) + abs(int(b.h)-int(a.h))
}

// AverageResolution returns the average resolution for all connected monitors.
func AverageResolution() (*Res, error) {
	monitors, err := Detect()
	if err != nil {
		return nil, err
	}
	var ws, hs uint
	for _, mon := range monitors {
		ws += mon.Width
		hs += mon.Height
	}
	ws /= uint(len(monitors))
	hs /= uint(len(monitors))
	return NewRes(ws, hs), nil
}

// Parses a string on the form "1234x1234"
func ParseSize(widthHeight string) (uint, uint, error) {
	fields := strings.SplitN(strings.ToLower(widthHeight), "x", 2)
	w, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, err
	}
	h, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, err
	}
	return uint(w), uint(h), nil
}

// FilenameToRes extracts width and height from a filename on the form: "asdf_123x123.xyz",
// or filenames that are just on the form "123x123.xyz".
func FilenameToRes(filename string) (*Res, error) {
	size := firstname(filepath.Base(filename))
	if strings.Contains(size, "_") {
		parts := strings.Split(size, "_")
		size = parts[len(parts)-1]
	}
	if !strings.Contains(size, "x") {
		return nil, errors.New("does not contain width x height: " + filename)
	}
	width, height, err := ParseSize(size)
	if err != nil {
		return nil, fmt.Errorf("does not contain width x height: %s: %s", filename, err)
	}
	return &Res{uint(width), uint(height)}, nil
}

// ExtractResolutions extracts Res structs from a slice of filenames
// All the filenames must be on the form *_WIDTHxHEIGHT.ext,
// where WIDTH and HEIGHT are numbers.
func ExtractResolutions(filenames []string) ([]*Res, error) {
	var resolutions []*Res
	for _, filename := range filenames {
		res, err := FilenameToRes(filename)
		if err != nil {
			return resolutions, err
		}
		resolutions = append(resolutions, res)
	}
	return resolutions, nil
}
