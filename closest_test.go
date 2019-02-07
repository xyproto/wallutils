package monitor

import (
	"fmt"
	"testing"
)

func TestClosest(t *testing.T) {
	filenames := []string{"hello_1024x768.jpg", "hello_1600x1200.jpg", "hello_320x200.jpg"}
	fmt.Println(filenames)

	resolutions, err := ExtractResolutions(filenames)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Resolutions extracted from slice of filenames:", resolutions)

	avgRes, err := AverageResolution()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Average resolution for all connected monitors:", avgRes)

	filename, err := Closest(filenames)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Filename closest to the average resolution:", filename)
}
