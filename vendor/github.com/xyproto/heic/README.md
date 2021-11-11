# HEIC

This is a fork of the Go code in [github.com/strukturag/libheif/go](https://github.com/strukturag/libheif/tree/master/go), with the intention of being able to load dynamic wallpapers in the `heic` format.

Code has been added that makes it possible to read the timing information from dynamic wallpapers in the `.heic` and/or HEIF format.

```go
...
metadataIDs := handle.MetadataIDs()
if len(metadataIDs) > 0 {
    metadataID := metadataIDs[0]
    timeTable, err := handle.ImageTimes(metadataID)
    // the mapping from image index to timestamp that contains the correct hour and minute are now in "timeTable"
    ...
}
```

Take a look at 

## General info

* Version: 1.0.0
* License: LGPL3
