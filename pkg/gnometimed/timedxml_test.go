package gnometimed

import (
	"fmt"
)

func ExampleParseXML() {
	gtw, err := ParseXML("testdata/example1.xml")
	if err != nil {
		panic(err)
	}
	fmt.Println(gtw.Config.StartTime.Year)
	fmt.Println(gtw.Config.Transitions[0].ToFilename)

	gtw, err = ParseXML("testdata/example2.xml")
	if err != nil {
		panic(err)
	}
	fmt.Println(gtw.Config.StartTime.Year)

	gtw, err = ParseXML("testdata/adwaita-timed.xml")
	if err != nil {
		panic(err)
	}
	fmt.Println(gtw.Config.StartTime.Year)

	gtw, err = ParseXML("testdata/generated.xml")
	if err != nil {
		panic(err)
	}
	fmt.Println(gtw.Config.StartTime.Year)

	// Output:
	// 2009
	// /usr/share/backgrounds/cosmos/comet.jpg
	// 0
	// 2011
	// 2018
}
