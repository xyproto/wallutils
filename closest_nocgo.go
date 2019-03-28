// +build !cgo

package wallutils

import "errors"

func ClosestByResolution(wallpapers []*Wallpaper) (*Wallpaper, error) {
	return nil, errors.New("not available when compiling without cgo")
}

func ClosestByResolutionInFilename(filenames []string) (string, error) {
	return "", errors.New("not available when compiling without cgo")
}

func Closest(filenames []string) (string, error) {
	return "", errors.New("not available when compiling without cgo")
}
