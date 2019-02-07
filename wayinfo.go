package monitor

// #cgo LDFLAGS: -lwayland-client
//#include "wayinfo.h"
import "C"
import (
	"errors"
	"log"
	"strconv"
	"strings"
)

func WaylandCanConnect() bool {
	return bool(C.WaylandRunning())
}

func WaylandInfo() (string, error) {
	if !WaylandCanConnect() {
		return "", errors.New("WaylandInfo(): not connected over Wayland")
	}
	infoString := C.GoString(C.WaylandInfoString())
	return infoString, nil
}

func WaylandMonitors(IDs, widths, heights, wDPIs, hDPIs *[]uint) error {
	if !WaylandCanConnect() {
		return errors.New("WaylandMonitors(): not connected over Wayland")
	}
	info, err := WaylandInfo()
	if err != nil {
		return err
	}

	var (
		counter     uint
		physCounter uint
		physW       uint
		physH       uint
		wDPI        uint
		hDPI        uint
	)

	// TODO: Write a C implementation instead of parsing the string output
	lines := strings.Split(info, "\n")

	// The physical width and height of the last encountered monitor, in millimeters
	for i, line := range lines {
		// Example output from Wayland Info:
		//   width: 1920 px, height: 1200 px, refresh: 60 Hz,
		//   flags: current
		if strings.Contains(line, "flags: current") {
			prevline := lines[i-1]
			fields := strings.Fields(prevline)
			if len(fields) > 4 {
				w, err := strconv.Atoi(fields[1])
				if err != nil {
					return err
				}
				h, err := strconv.Atoi(fields[4])
				if err != nil {
					return err
				}
				*IDs = append(*IDs, counter)
				*widths = append(*widths, uint(w))
				*heights = append(*heights, uint(h))
				if physW == 0 || physH == 0 {
					wDPI = 96 // default DPI value, if no physical size is given
					hDPI = 96 // default DPI value, if no physical size is given
					log.Println("WARN: No physical monitor size detected!")
					//return errors.New("no physical monitor size detected")
				}
				if physW > 0 && physH > 0 {
					// Calculate DPI, from the monitor size (in mm) and the pixel size
					wDPI = uint(float64(w) / (float64(physW) / 25.4))
					hDPI = uint(float64(h) / (float64(physH) / 25.4))
				}
				*wDPIs = append(*wDPIs, wDPI)
				*hDPIs = append(*hDPIs, hDPI)
				counter++
			}
		} else if strings.Contains(line, "physical_width:") {
			// Example output from Wayland Info:
			//   physical_width: 518 mm, physical_height: 324 mm,
			fields := strings.Fields(line)
			if len(fields) > 5 {
				w, err := strconv.Atoi(fields[1])
				if err != nil {
					return err
				}
				h, err := strconv.Atoi(fields[4])
				if err != nil {
					return err
				}
				physW = uint(w)
				physH = uint(h)
				physCounter++
			}
		}
	}
	if 0 < physCounter && physCounter < counter {
		log.Println("WARN: Some monitors contains a physical size, but not all of them")
		//return errors.New("some monitors contains a physical size, but not all of them")
	}
	return nil
}
