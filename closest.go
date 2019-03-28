// +build cgo

package wallutils

// ClosestByResolution returns a wallpaper that is closest to the average
// monitor resolution. If several wallpapers matches, a random one is returned.
// The idea is that a slice of wallpapers in a wallpaper collection with several
// available resolutions is given as input, and a suitable wallpaper is returned.
func ClosestByResolution(wallpapers []*Wallpaper) (*Wallpaper, error) {
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

// ClosestByResolutionInFilename takes a list of filenames on the form
// "*_WIDTHxHEIGHT.ext", where WIDTH and HEIGHT are numbers.
// The filename that is closest to the average monitor resolution is returned.
// Any filenames not following the pattern will cause an error being returned.
func ClosestByResolutionInFilename(filenames []string) (string, error) {
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

// Closest does the same as ClosestByResolutionInfilename.
// It is provided for backwards compatibility.
func Closest(filenames []string) (string, error) {
	return ClosestByResolutionInFilename(filenames)
}
