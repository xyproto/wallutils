# Monitor [![Build Status](https://travis-ci.org/xyproto/monitor.svg?branch=master)](https://travis-ci.org/xyproto/monitor) [![GoDoc](https://godoc.org/github.com/xyproto/monitor?status.svg)](http://godoc.org/github.com/xyproto/monitor) [![License](http://img.shields.io/badge/license-MIT-green.svg?style=flat)](https://raw.githubusercontent.com/xyproto/monitor/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/monitor)](https://goreportcard.com/report/github.com/xyproto/monitor)

Detect monitor resolutions and set the desktop wallpaper, for any windowmanager (please file an issue if your windowmanager is not supported yet). Also supports timed wallpapers.

## Highlighted functions

* The `monitor.Detect` function can return a `[]Monitor` slice in Go, with information about all connected monitors.
* The `monitor.SetTimedWallpaper` function can be used for launching an event loop that handles GNOME timed wallpapers.

## Included utilities

  * `getdpi`, for retrieving the average DPI, for all monitors.
  * `lscollections`, for listing installed wallpaper collections.
  * `timedinfo`, for showing more information about installed GNOME timed wallpapers.
  * `lsmon` lists the connected monitors and resolutions.
  * `lstimed` for listing installed GNOME timed wallpapers.
  * `lswallpaper`, for listing all installed wallpapers.
  * `setcollection`, for setting a suitable (in terms of resolution) wallpaper from a wallpaper collection.
  * `setrandom`, for setting a random wallpaper.
  * `settimed`, for setting GNOME timed wallpapers (will continue to run, to handle time events).
  * `setwallpaper` can be used for setting a wallpaper (works both for X11 and Wayland).
  * `wayinfo` shows detailed information about the connected monitors, via Wayland.
  * `xinfo` shows detailed information about the connected monitors, via X11.

## Example use of the `lsmon` utility

```sh
$ lsmon
0: 1920x1200
1: 1920x1200
2: 1920x1200
```

## Building and installing utilities

Using make, for building and installing all included utilities:

    make
    make install

Using Go 1.11 or later, for a single utility:

    go get -u github.com/xyproto/monitor/cmd/setwallpaper

On Ubuntu, from a fresh installation:

    sudo apt get update
    sudo apt get install libxcursor-dev libxmu-dev libx11-dev git golang-go
    go get -u github.com/xyproto/monitor/cmd/setwallpaper
    cd ~/go/src/github.com/xyproto/monitor

Manually:

    # clone the repository
    git clone https://github.com/xyproto/monitor

    # build and install the setmonitor command
    cd monitor/cmd/setmonitor
    go build
    install -Dm755 setmonitor /usr/bin/setmonitor

## Example use of `setwallpaper`

    setwallpaper /path/to/background/image.png

## Example use of the Go package

### Retrieve monitor resolution(s)

~~~go
package main

import (
	"fmt"
	"os"

	"github.com/xyproto/monitor"
)

func main() {
	// Retrieve a slice of Monitor structs, or exit with an error
	monitors, err := monitor.Detect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	// For every monitor, output the ID, width and height
	for _, monitor := range monitors {
		fmt.Printf("%d: %dx%d\n", monitor.ID, monitor.Width, monitor.Height)
	}
}
~~~

### Change the wallpaper

```go
fmt.Println("Setting background image to: " + imageFilename)
if err := monitor.SetWallpaper(imageFilename); err != nil {
	return err
}
```

## Build requirements

* Go 1.11 or later.
* A working C compiler (tested with GCC 8.2.1)
* `libwayland-client.so` and header files available, for Wayland support.
* `libX11.so` and header files available, for X11 support.

## General info

* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
* License: MIT
* Version: 3.2.0
