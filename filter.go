package monitor

// FilterWallpapers will filter out wallpapers that both match with the collection name, and are also marked as part of a collection
func FilterWallpapers(collectionName string, wallpapers []*Wallpaper) []*Wallpaper {
	var collection []*Wallpaper
	for _, wp := range wallpapers {
		if wp.PartOfCollection && wp.CollectionName == collectionName {
			collection = append(collection, wp)
		}
	}
	return collection
}

// FilterSimpleTimedWallpapers will filter out simpleTimed timed wallpapers that match with the collection name
func FilterSimpleTimedWallpapers(collectionName string, simpleTimedWallpapers []*SimpleTimedWallpaper) []*SimpleTimedWallpaper {
	var collection []*SimpleTimedWallpaper
	for _, stw := range simpleTimedWallpapers {
		if stw.Name == collectionName {
			collection = append(collection, stw)
		}
	}
	return collection
}

// FilterGnomeWallpapers will filter out gnome timed wallpapers that match with the collection name
func FilterGnomeWallpapers(collectionName string, gnomeWallpapers []*GnomeWallpaper) []*GnomeWallpaper {
	var collection []*GnomeWallpaper
	for _, gw := range gnomeWallpapers {
		if gw.CollectionName == collectionName {
			collection = append(collection, gw)
		}
	}
	return collection
}
