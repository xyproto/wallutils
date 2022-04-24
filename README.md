# Wallutils [![Build](https://github.com/xyproto/wallutils/actions/workflows/build.yml/badge.svg)](https://github.com/xyproto/wallutils/actions/workflows/build.yml) [![GoDoc](https://godoc.org/github.com/xyproto/wallutils?status.svg)](http://godoc.org/github.com/xyproto/wallutils) [![License](http://img.shields.io/badge/license-BSD-green.svg?style=flat)](https://raw.githubusercontent.com/xyproto/wallutils/main/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/wallutils)](https://goreportcard.com/report/github.com/xyproto/wallutils)

* Detect monitor resolutions and set the desktop wallpaper, for any window manager (please file an issue if your window manager is not supported yet).
* Supports GNOME timed wallpapers, and includes a utility that can run an event loop for changing them (also supports cross fading).
* Introduces a new file format for timed wallpapers: The **Simple Timed Wallpaper** format: [Markdown](https://github.com/xyproto/wallutils/blob/main/pkg/simpletimed/stw-1.0.0.md) | [PDF](https://raw.githubusercontent.com/xyproto/wallutils/main/pkg/simpletimed/stw-1.0.0.pdf).
* GNOME timed wallpapers can be converted to the Simple Timed Wallpaper format with the `xml2stw` utility.
* macOS dynamic wallpapers (in the HEIF format with the `.heic` extension) can be installed with `heic-install` and used with `lstimed` and `settimed`. This extracts the metadata with `heic2stw` (only timing information, not the azimuth and elevation for the sun, yet) and extracts the images with `convert` that comes with ImageMagick.

[![Packaging status](https://repology.org/badge/vertical-allrepos/wallutils.svg)](https://repology.org/project/wallutils/versions)

## Timed Wallpapers

The [Mojave timed wallpaper](https://github.com/japamax/gnome-mojave-timed-wallpaper) and other timed wallpapers can be set with the `settimed` command, and will cross fade from image to image as the day progresses:

<img alt="Dynamic wallpaper example" src="https://i.redd.it/z5zx32pe3l311.gif" width=320>

## Included utilities

  * `getdpi`, for retrieving the average DPI, for all monitors (use `-b` to see the DPI both horizontally and vertically).
  * `lscollection`, for listing installed wallpaper collections (use `-l` for also listing paths and collection names).
  * `timedinfo`, for showing more information about installed timed wallpapers.
  * `lsmon` lists the connected monitors and resolutions (use `-d` for also listing DPI).
  * `lstimed` for listing installed timed wallpapers (use `-l` for also listing paths).
  * `lswallpaper`, for listing all installed wallpapers (use `-l` and `-s` to list more information).
  * `setcollection`, for setting a suitable (in terms of resolution) wallpaper from a wallpaper collection.
  * `setrandom`, for setting a random wallpaper.
  * `settimed`, for setting timed wallpapers (will continue to run, to handle time events). (This utility has recently been refactored and needs more testing).
  * `setwallpaper` can be used for setting a wallpaper (works both over X and the Wayland protocol).
  * `wayinfo` shows detailed information about the connected monitors, via Wayland.
  * `xinfo` shows detailed information about the current X setup.
  * `xml2stw` for converting GNOME timed wallpapers to the Simple Timed Wallpaper format.
  * `heic2stw` for extracting the timing information from macOS dynamic wallpapers (`.heic` files) to the Simple Timed Wallpaper format.

## Included scripts

  * `heic-install` for installing a macOS dynamic wallpaper to `/usr/share/backgrounds` using both ImageMagick `convert` and `heic2stw`.

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
    sudo dnf install git golang imagemagick libXcursor-devel libXmu-devel xorg-x11-xbitmaps libheif-devel
    git clone https://github.com/xyproto/wallutils
    cd wallutils
    make
    sudo make PREFIX=/usr/local install

### Debian 11

With sudo, or as root, install the required dependencies:

    apt install git golang imagemagick libx11-dev libxcursor-dev libxmu-dev libwayland-dev libxpm-dev xbitmaps libxmu-headers libheif-dev make

As a user, clone the repository and build the utilities:

    git clone https://github.com/xyproto/wallutils
    cd wallutils
    make

Then with sudo, or as root, again, install wallutils:

    sudo make PREFIX=/usr/local install

## Installing a single utility

Using Go 1.17 or later, installing the `getdpi` utility:

    go install github.com/xyproto/wallutils/cmd/getdpi@latest

## Wayland or X only

The packages related to X can be removed after building if you only wish to keep the Wayland-related functionality. And the same for X.

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

* Go 1.11 or later. 1.17 or later is recommended.
* A working C compiler (tested with GCC 8.2.1).
* Header files for Wayland and X.
* `libheif` for `heic2stw` and `heic-install`.

## Runtime requirements

* `libwayland-client.so`, for Wayland support.
* `libX11.so`, for X support.
* `libheif.so` for HEIF/`.heic` image support.

It is also possible to build with `make static`, to only build the utilities that does not depend on any of the above `.so` files, as statically compiled ELF executables.

## Wallpaper collections

The XML format from GNOME for specifying wallpaper **collections** is not yet supported (and I'm not sure if it's needed). Creating a directory with images where the filename of the images specify the resolution (like `wallpaper_5639x3561.jpg`) is enough for `lscollection` to recognize it as a collection (if the directory is placed in `/usr/share/backgrounds` or `/usr/share/wallpapers`).

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
* License: BSD-3
* Version: 5.10.0
