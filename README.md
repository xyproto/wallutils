# Wallutils [![Build Status](https://travis-ci.org/xyproto/wallutils.svg?branch=master)](https://travis-ci.org/xyproto/wallutils) [![GoDoc](https://godoc.org/github.com/xyproto/wallutils?status.svg)](http://godoc.org/github.com/xyproto/wallutils) [![License](http://img.shields.io/badge/license-MIT-green.svg?style=flat)](https://raw.githubusercontent.com/xyproto/wallutils/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/wallutils)](https://goreportcard.com/report/github.com/xyproto/wallutils)

* Detect monitor resolutions and set the desktop wallpaper, for any window manager (please file an issue if your window manager is not supported yet).
* Supports GNOME timed wallpapers, and includes a utility that can run an event loop for changing them (also supports cross fading).
* Introduces a new file format for timed wallpapers: The **Simple Timed Wallpaper** format: [Web](https://github.com/xyproto/simpletimed/#specification) | [Markdown](https://github.com/xyproto/simpletimed/blob/master/stw-1.0.0.md) | [PDF](https://github.com/xyproto/simpletimed/raw/master/stw-1.0.0.pdf)

## Timed Wallpapers

The [Mojave timed wallpaper](https://github.com/japamax/gnome-mojave-timed-wallpaper) and other timed wallpapers can be set with the `settimed` command, and will cross fade from image to image as the day progresses:

![](https://i.redd.it/z5zx32pe3l311.gif)]

## Included utilities

  * `getdpi`, for retrieving the average DPI, for all monitors (use `-b` to see the DPI both horizontally and vertically).
  * `lscollection`, for listing installed wallpaper collections (use `-l` for also listing paths and collection names).
  * `timedinfo`, for showing more information about installed timed wallpapers.
  * `lsmon` lists the connected monitors and resolutions (use `-d` for also listing DPI).
  * `lstimed` for listing installed timed wallpapers (use `-l` for also listing paths).
  * `lswallpaper`, for listing all installed wallpapers (use `-l` and `-s` to list more information).
  * `setcollection`, for setting a suitable (in terms of resolution) wallpaper from a wallpaper collection.
  * `setrandom`, for setting a random wallpaper.
  * `settimed`, for setting timed wallpapers (will continue to run, to handle time events).
  * `setwallpaper` can be used for setting a wallpaper (works both over X and the Wayland protocol).
  * `wayinfo` shows detailed information about the connected monitors, via Wayland.
  * `xinfo` shows detailed information about the current X setup.
  * `xml2stw` for converting GNOME timed wallpapers to the Simple Timed Wallpaper format.

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

Using Go 1.12 or later, for a single utility:

    go get -u github.com/xyproto/wallutils/cmd/settimed

On Arch Linux:

Install `wallutils` from AUR, or:

    sudo pacman -Syu git go libxcursor libxmu wayland xbitmaps xorgproto
    git clone https://github.com/xyproto/wallutils
    cd wallutils
    make
    sudo make install

On Fedora:

    sudo dnf update
    sudo dnf install xorg-x11-xbitmaps libXcursor-devel libXmu-devel
    git clone https://github.com/xyproto/wallutils
    cd wallutils
    make
    sudo make install

On Ubuntu:

    sudo apt get update
    sudo apt get install libxcursor-dev libxmu-dev libx11-dev git golang-go
    git clone https://github.com/xyproto/wallutils
    cd wallutils
    make
    sudo make install

## Wayland or X only

The packages related to X can be removed after building if only wish to keep the Wayland-related functionality. And likewise for X.

## Example use of `settimed`

    settimed mojave-timed

## Example use of `setwallpaper`

    setwallpaper /path/to/background/image.png

## Example use of `setrandom`

    setrandom /usr/share/pixmaps

## Example use of the Go package

### Retrieve monitor resolution(s)

~~~go
package main

import (
	"fmt"
	"os"

	"github.com/xyproto/wallutils"
)

func main() {
	// Retrieve a slice of Monitor structs, or exit with an error
	monitors, err := wallutils.Monitors()
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
if err := wallutils.SetWallpaper(imageFilename); err != nil {
	return err
}
```

## Build requirements

* Go 1.12 or later.
* A working C compiler (tested with GCC 8.2.1).
* Header files for Wayland and X.

## Runtime requirements

* `libwayland-client.so`, for Wayland support.
* `libX11.so`, for X support.

## General info

* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
* License: MIT
* Version: 5.4.0
