package monitor

// #cgo LDFLAGS: -lX11
// #include "xinfo.h"
import "C"
import (
	"errors"
	"strconv"
	"strings"
)

func XCanConnect() bool {
	return bool(C.X11Running())
}

func XInfo() (string, error) {
	if !XCanConnect() {
		return "", errors.New("XInfo(): not connected over X11")
	}
	infoString := C.GoString(C.X11InfoString())
	return infoString, nil
}

// Monitor enumerates the monitors and returns a slice of structs,
// including the resolution.
func XMonitors(IDs, widths, heights, wDPIs, hDPIs *[]uint) error {
	if !XCanConnect() {
		return errors.New("XMonitors(): not connected over X11")
	}
	info, err := XInfo()
	if err != nil {
		return err
	}
	var counter uint
	// TODO: Write a C implementation instead of parsing the info string
	for _, line := range strings.Split(info, "\n") {
		if strings.Contains(line, "dimensions:") {
			fields := strings.Fields(line)
			if len(fields) > 2 && strings.Contains(fields[1], "x") {
				resFields := strings.SplitN(fields[1], "x", 2)
				w, err := strconv.Atoi(resFields[0])
				if err != nil {
					return err
				}
				h, err := strconv.Atoi(resFields[1])
				if err != nil {
					return err
				}
				*IDs = append(*IDs, counter)
				*widths = append(*widths, uint(w))
				*heights = append(*heights, uint(h))
				counter++
			}
		} else if strings.Contains(line, "resolution:") {
			fields := strings.Fields(line)
			if len(fields) > 2 && strings.Contains(fields[1], "x") {
				dpiFields := strings.SplitN(fields[1], "x", 2)
				wDPI, err := strconv.Atoi(dpiFields[0])
				if err != nil {
					return err
				}
				hDPI, err := strconv.Atoi(dpiFields[1])
				if err != nil {
					return err
				}
				*wDPIs = append(*wDPIs, uint(wDPI))
				*hDPIs = append(*hDPIs, uint(hDPI))
			}
		}
	}
	return nil
}
