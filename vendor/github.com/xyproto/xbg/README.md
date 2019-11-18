# xbg [![Build Status](https://travis-ci.com/xyproto/xbg.svg?branch=master)](https://travis-ci.com/xyproto/xbg) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/xbg)](https://goreportcard.com/report/github.com/xyproto/xbg) [![GoDoc](https://godoc.org/github.com/xyproto/xbg?status.svg)](https://godoc.org/github.com/xyproto/xbg)

Go module for setting the background image under X11.

* Can be used under window managers like i3 and AwesomeWM.
* The C code is based on code from [bgs](https://github.com/Gottox/bgs) (also licensed under MIT/X, credited in the LICENSE file).
* Intended to be used by [wallutils](https://github.com/xyproto/wallutils).

## Example

A test-utility is included in `cmd/grumpybg/`.

## Plans

- [ ] Rewrite the C code using a Go module that can use the X11 protocol directly.
- [ ] Support multiple monitors in `x11.go`.
- [ ] Support monitor rotation in `x11.go`.

## General info

* Version: 0.0.5
* License: MIT
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;

`bgs.h` is a modified version of code that is licensed under the MIT/X license, `copyright © 2007-2008 Enno Gottox Boland <gottox at s01 dot de>`.
