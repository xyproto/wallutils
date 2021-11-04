package gnometimed

import (
	"fmt"
	"testing"
)

func TestConvert(t *testing.T) {
	gtw, err := ParseXML("testdata/generated.xml")
	if err != nil {
		panic(err)
	}
	fmt.Println(gtw.Config.StartTime.Year)
	stw, err := GnomeToSimple(gtw)
	if err != nil {
		panic(err)
	}
	fmt.Println(stw)
}
