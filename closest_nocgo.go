//go:build !cgo
// +build !cgo

package wallutils

import "errors"

func ClosestByResolution(wallpapers []*Wallpaper) (*Wallpaper, error) {
	return nil, errors.New("the function ClosestByResolution is not available when compiling without cgo")
}

func ClosestByResolutionInFilename(filenames []string) (string, error) {
	return "", errors.New("the function ClosestByResolutionInFilename is not available when compiling without cgo")
}

func Closest(filenames []string) (string, error) {
	return "", errors.New("the function Closest is not available when compiling without cgo")
}
