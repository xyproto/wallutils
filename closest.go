package wallutils

import (
	"os"
)

// Exists checks if the given filename exists in the current directory
// (or if an absolute path exists)
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// Closest takes a list of filenames on the form "*_WIDTHxHEIGHT.ext".
// WIDTH and HEIGHT are numbers. Closest returns the filename that is closest
// to the average monitor resolution. Any filenames not following the pattern
// will result in an error being returned.
func Closest(filenames []string) (string, error) {
	avgRes, err := AverageResolution()
	if err != nil {
		return "", err
	}
	// map: (distance to average resolution) => (filename)
	d := make(map[int]string)
	var dist int
	var minDist int
	var minDistSet bool
	for _, filename := range filenames {
		res, err := FilenameToRes(filename)
		if err != nil {
			return "", err
		}
		dist = Distance(avgRes, res)
		if dist < minDist || !minDistSet {
			minDist = dist
			minDistSet = true
		}
		//fmt.Printf("FILENAME %s HAS DISTANCE %d TO AVERAGE RESOLUTION %s\n", filename, dist, avgRes)
		d[dist] = filename
	}
	// ok, have a map, now find the filename of the smallest distance
	return d[minDist], nil
}

// GoodFit returns the image file with a resolution that is closest to the
// current everage monitor resolution.
func GoodFit(wallpapers []*Wallpaper) (*Wallpaper, error) {
	avgRes, err := AverageResolution()
	if err != nil {
		return nil, err
	}
	// map: "distance to average resolution" => wallpaper
	d := make(map[int](*Wallpaper))
	var dist int
	var minDist int
	var minDistSet bool
	for _, wp := range wallpapers {
		res := wp.Res()
		dist = Distance(avgRes, res)
		if dist < minDist || !minDistSet {
			minDist = dist
			minDistSet = true
		}
		d[dist] = wp
	}
	// ok, have a map, now find the filename of the smallest distance
	return d[minDist], nil
}
