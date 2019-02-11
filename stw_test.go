package monitor

import (
	"fmt"
)

func ExampleParseSTW() {
	stw, err := ParseSTW("testdata/adwaita-timed2.stw")
	if err != nil {
		panic(err)
	}
	fmt.Println(stw.Name)

	// Output:
	// adwaita-timed
}
