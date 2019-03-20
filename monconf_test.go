package wallutils

import (
	"testing"
)

func TestParseMonitorConfiguration(t *testing.T) {
	mc, err := NewMonitorConfiguration()
	if err != nil {
		// Ignore this test if ~/.config/monitor.xml does not exist
		return
	}
	mc.Overlapping()
}
