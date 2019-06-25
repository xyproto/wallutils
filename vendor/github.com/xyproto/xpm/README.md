# xpm

Encode images to the X PixMap (XPM3) image format.

The resulting images are smaller than the one from GIMP, since the question mark character is also used, while at the same time avoiding double question marks, which could result in a trigraph (like `??=`, which has special meaning in C).

Includes a `png2xpm` utility.

* Version: 1.2.1
* License: MIT
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
