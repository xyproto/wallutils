package wallutils

import (
	"fmt"
	"strings"
	"testing"
)

func TestXrandrOverlap(t *testing.T) {
	if overlap, reslines := XRandrOverlap(); overlap {
		fmt.Println("XRandr: monitor configurations are overlapping!")
		fmt.Println(strings.Join(reslines, "\n"))
	} else {
		fmt.Println("No overlapping monitor configurations.")
	}
}
