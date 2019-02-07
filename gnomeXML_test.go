package monitor

import (
	"fmt"
)

func ExampleParse() {
	b, err := Parse("testdata/example1.xml")
	if err != nil {
		panic(err)
	}
	fmt.Println(b.StartTime.Year)
	fmt.Println(b.Transitions[0].ToFilename)

	// ---

	b, err = Parse("testdata/example2.xml")
	if err != nil {
		panic(err)
	}
	fmt.Println(b.StartTime.Year)

	// ---

	b, err = Parse("testdata/adwaita-timed.xml")
	if err != nil {
		panic(err)
	}
	fmt.Println(b.StartTime.Year)

	// Output:
	// 2009
	// /usr/share/backgrounds/cosmos/comet.jpg
	// 0
	// 2011
}
