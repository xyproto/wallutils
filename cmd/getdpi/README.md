# getdpi

Tool for retrieving the average DPI across all monitors, regardless of if X or Wayland is in use.

## Building getdpi

    go build

## Usage

Retreive the average DPI as a pair of numbers (ie. `96x96`):

    getdpi -b

Retreive the average DPI as a single number (ie. `96`):

    getdpi

Version information:

    getdpi --version

## Listing DPI per monitor

The `lsmon` utility in the [monitor](https://github.com/xyproto/monitor) package supports listing monitor resolutions and DPI with the `-dpi` flag.

## General Info

* Version: 1.2.0
* License: MIT
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
