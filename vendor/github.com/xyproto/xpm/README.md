# xpm [![Build Status](https://travis-ci.com/xyproto/xpm.svg?branch=master)](https://travis-ci.com/xyproto/xpm) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/xpm)](https://goreportcard.com/report/github.com/xyproto/xpm) [![GoDoc](https://godoc.org/github.com/xyproto/xpm?status.svg)](https://godoc.org/github.com/xyproto/xpm)

Encode images to the X PixMap (XPM3) image format.

The resulting images are smaller than the one from GIMP, since the question mark character is also used, while at the same time avoiding double question marks, which could result in a trigraph (like `??=`, which has special meaning in C).

Includes a `png2xpm` utility.

* Version: 2.1.0
* License: MIT
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
