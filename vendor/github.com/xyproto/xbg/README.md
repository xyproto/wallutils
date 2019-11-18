# xbg [![Build Status](https://travis-ci.com/xyproto/xbg.svg?branch=master)](https://travis-ci.com/xyproto/xbg) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/xbg)](https://goreportcard.com/report/github.com/xyproto/xbg) [![GoDoc](https://godoc.org/github.com/xyproto/xbg?status.svg)](https://godoc.org/github.com/xyproto/xbg)

Go module for setting the background image under X (Xorg/X11).

* Can be used together with window managers like `i3` and `AwesomeWM`.
* Intended to be used by [wallutils](https://github.com/xyproto/wallutils).
* Based on code from [bgs](https://github.com/Gottox/bgs) by Enno Gottox Boland, which uses Imlib2.
* By using `New()` and `.Release()`, the X11 struct is thread-safe. The `SetWallpaper` function uses these.

## Example

A test-utility is included in `cmd/grumpybg/`.

## Plans

- [ ] Rewrite the C code using a Go module that can use the X11 protocol directly, then make it thread-safe.
- [ ] Support multiple monitors in `x11.go`.
- [ ] Support monitor rotation in `x11.go`.

## General info

* Version: 0.2.0
* License: MIT
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;

`bgs.h` is a modified version of `bgs.c` that is licensed under the MIT/X license and is copyright © 2007-2008 Enno Gottox Boland <gottox at s01 dot de>.
