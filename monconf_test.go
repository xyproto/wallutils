package wallutils

import (
	"fmt"
	"testing"
)

func TestParseMonitorConfiguration(t *testing.T) {
	mc, err := ParseMonitors()
	if err != nil {
		panic(err)
	}
	for _, conf := range mc.Configurations {
		for _, output := range conf.Outputs {
			fmt.Println("OUTPUT", output)
		}
	}
	fmt.Println("NO OUTPUT")
}
