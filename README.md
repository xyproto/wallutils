# Wallutils [![Build Status](https://travis-ci.com/xyproto/wallutils.svg?branch=master)](https://travis-ci.com/xyproto/wallutils) [![GoDoc](https://godoc.org/github.com/xyproto/wallutils?status.svg)](http://godoc.org/github.com/xyproto/wallutils) [![License](http://img.shields.io/badge/license-MIT-green.svg?style=flat)](https://raw.githubusercontent.com/xyproto/wallutils/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/wallutils)](https://goreportcard.com/report/github.com/xyproto/wallutils)

* Detect monitor resolutions and set the desktop wallpaper, for any window manager (please file an issue if your window manager is not supported yet).
* Supports GNOME timed wallpapers, and includes a utility that can run an event loop for changing them (also supports cross fading).
* Introduces a new file format for timed wallpapers: The **Simple Timed Wallpaper** format: [Web](https://github.com/xyproto/simpletimed/#specification) | [Markdown](https://github.com/xyproto/simpletimed/blob/master/stw-1.0.0.md) | [PDF](https://github.com/xyproto/simpletimed/raw/master/stw-1.0.0.pdf)

[![Packaging status](https://repology.org/badge/vertical-allrepos/wallutils.svg)](https://repology.org/project/wallutils/versions)

## Timed Wallpapers

The [Mojave timed wallpaper](https://github.com/japamax/gnome-mojave-timed-wallpaper) and other timed wallpapers can be set with the `settimed` command, and will cross fade from image to image as the day progresses:

![](https://i.redd.it/z5zx32pe3l311.gif)

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

## Installing `wallutils`

### Arch Linux

    pacman -S wallutils

### Fedora

Until an official package is available:

    sudo dnf update
    sudo dnf install git golang libXcursor-devel libXmu-devel xorg-x11-xbitmaps
    git clone https://github.com/xyproto/wallutils
    cd wallutils
    make
    sudo make PREFIX=/usr/local install

### Ubuntu

Until an official package is available:

Go 1.11 or later is required, [here's an easy way to install Go 1.12](https://github.com/golang/go/wiki/Ubuntu):

    sudo add-apt-repository ppa:longsleep/golang-backports
    sudo apt-get update
    sudo apt-get install golang-go

Then install the required dependencies, clone the repository and install wallutils:

    sudo apt get install git libx11-dev libxcursor-dev libxmu-dev libwayland-dev libxpm-dev xbitmaps libxmu-headers
    git clone https://github.com/xyproto/wallutils
    cd wallutils
    make
    sudo make PREFIX=/usr/local install

## Installing a single utility

Using Go 1.11 or later, installing the `settimed` utility:

    go get -u github.com/xyproto/wallutils/cmd/settimed

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

* Go 1.11 or later.
* A working C compiler (tested with GCC 8.2.1).
* Header files for Wayland and X.

## Runtime requirements

* `libwayland-client.so`, for Wayland support.
* `libX11.so`, for X support.

It is also possible to build with `make static`, to only build the utilities that does not depend on any of the above `.so` files, as statically compiled ELF executables.

## Wallpaper collections

The XML format from GNOME for specifying wallpaper collections is not yet supported (and I'm not sure if it's needed). Creating a directory with images where the filename of the images specify the resolution (like `wallpaper_5639x3561.jpg`) is enough for `lscollection` to recognize it as a collection (if the directory is placed in `/usr/share/backgrounds` or `/usr/share/wallpapers`).

## Refreshing the wallpaper after waking from sleep

Send the `USR1` signal to the `settimed` process:

    pkill settimed -USR1

This should refresh the wallpaper.

## A note about i3

* When using wallutils together with `i3`, it works best with also having `feh` and `imlib2` installed.

## Setting a wallpaper per monitor

* Setting a wallpaper per monitor is not supported, yet. Currently, a wallpaper is set for all monitors.

## General info

* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
* License: MIT
* Version: 5.9.0
