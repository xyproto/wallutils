package wallutils

import (
	"fmt"
	"log"
	"testing"
)

func TestXMonitors(t *testing.T) {
	if !XCanConnect() {
		log.Println("Could not connect over X, skipping test")
		return
	}
	info, err := XInfo()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(info)
	IDs, widths, heights, wDPIs, hDPIs := []uint{}, []uint{}, []uint{}, []uint{}, []uint{}
	err = XMonitors(&IDs, &widths, &heights, &wDPIs, &hDPIs)
	if err != nil {
		t.Error(err)
	}
	for i, ID := range IDs {
		fmt.Printf("monitor %d: %dx%d (DPI: %dx%d)\n", ID, widths[i], heights[i], wDPIs[i], hDPIs[i])
	}
}
