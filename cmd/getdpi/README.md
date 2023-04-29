# getdpi

Tool for retrieving the average DPI across all monitors, regardless of if X or Wayland is in use.

## Building getdpi

    go build

## Usage

Retrieve the average DPI as a pair of numbers (ie. `96x96`):

    getdpi -b

Retrieve the average DPI as a single number (ie. `96`):

    getdpi

Version information:

    getdpi --version

## Listing DPI per monitor

The `lsmon` utility supports listing monitor resolutions and DPI with the `-dpi` flag.
