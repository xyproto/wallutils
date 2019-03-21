package wallutils

import (
	"fmt"
	"testing"
)

func TestParseMonitorConfiguration(t *testing.T) {
	mc, err := NewMonitorConfiguration()
	if err != nil {
		// Ignore this test if ~/.config/monitor.xml does not exist
		return
	}
	if !mc.Overlapping() {
		fmt.Println("No overlapping monitor configurations in ~/.config/monitors.xml")
	} else {
		fmt.Println("There are overlapping monitor configurations in ~/.config/monitors.xml")
	}
}
