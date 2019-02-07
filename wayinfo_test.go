package monitor

import (
	"fmt"
	"log"
	"testing"
)

func TestWaylandMonitors(t *testing.T) {
	if !WaylandCanConnect() {
		log.Println("Could not connect over Wayland, skipping test")
		return
	}

	info, err := WaylandInfo()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(info)
	IDs, widths, heights, wDPIs, hDPIs := []uint{}, []uint{}, []uint{}, []uint{}, []uint{}
	err = WaylandMonitors(&IDs, &widths, &heights, &wDPIs, &hDPIs)
	if err != nil {
		t.Error(err)
	}
	if len(IDs) == 0 {
		t.Fatal("No monitors?")
	}
	for i, ID := range IDs {
		fmt.Printf("monitor %d: %dx%d (DPI: %dx%d)\n", ID, widths[i], heights[i], wDPIs[i], hDPIs[i])
	}
}
