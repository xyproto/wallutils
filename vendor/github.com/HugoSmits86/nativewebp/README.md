[![Codecov Coverage](https://codecov.io/gh/HugoSmits86/nativewebp/branch/main/graph/badge.svg)](https://codecov.io/gh/HugoSmits86/nativewebp)
[![Go Reference](https://pkg.go.dev/badge/github.com/HugoSmits86/nativewebp.svg)](https://pkg.go.dev/github.com/HugoSmits86/nativewebp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# Native WebP for Go

This is a native WebP encoder written entirely in Go, with **no dependencies on libwebp** or other external libraries. Designed for performance and efficiency, this encoder generates smaller files than the standard Go PNG encoder and is approximately **50% faster** in execution.

Currently, the encoder supports only WebP lossless images (VP8L).

## Decoding Support

We provide WebP decoding through a wrapper around `golang.org/x/image/webp`, with an additional `DecodeIgnoreAlphaFlag` function to handle VP8X images where the alpha flag causes decoding issues.
## Benchmark

We conducted a quick benchmark to showcase file size reduction and encoding performance. Using an image from Google’s WebP Lossless and Alpha Gallery, we compared the results of our nativewebp encoder with the standard PNG encoder. <br/><br/>
For the PNG encoder, we used `png.BestCompression`. Likewise, nativewebp was configured with `nativewebp.BestCompression` so both encoders were benchmarked using their maximum compression settings.
<br/><br/>

<table align="center">
  <tr>
    <th></th>
    <th></th>
    <th>PNG encoder</th>
    <th>nativeWebP encoder</th>
    <th>reduction</th>
  </tr>
  <tr>
    <td rowspan="2" height="110px"><p align="center"><img src="https://www.gstatic.com/webp/gallery3/1.png" height="100px"></p></td>
    <td>file size</td>
    <td>120 kb</td>
    <td>95 kb</td>
    <td>21% smaller</td>
  </tr>
  <tr>
    <td>encoding time</td>
    <td>42945049 ns/op</td>
    <td>35413726 ns/op</td>
    <td>17% faster</td>
  </tr>
  <tr>
    <td rowspan="2" height="110px"><p align="center"><img src="https://www.gstatic.com/webp/gallery3/2.png" height="100px"></p></td>
    <td>file size</td>
    <td>46 kb</td>
    <td>35 kb</td>
    <td>24% smaller</td>
  </tr>
  <tr>
    <td>encoding time</td>
    <td>98509399 ns/op</td>
    <td>42626779 ns/op</td>
    <td>57% faster</td>
  </tr>
  <tr>
    <td rowspan="2" height="110px"><p align="center"><img src="https://www.gstatic.com/webp/gallery3/3.png" height="100px"></p></td>
    <td>file size</td>
    <td>236 kb</td>
    <td>190 kb</td>
    <td>19% smaller</td>
  </tr>
  <tr>
    <td>encoding time</td>
    <td>178205535 ns/op</td>
    <td>96800750 ns/op</td>
    <td>46% faster</td>
  </tr>
  <tr>
    <td rowspan="2" height="110px"><p align="center"><img src="https://www.gstatic.com/webp/gallery3/4.png" height="60px"></p></td>
    <td>file size</td>
    <td>53 kb</td>
    <td>39 kb</td>
    <td>26% smaller</td>
  </tr>
  <tr>
    <td>encoding time</td>
    <td>29088555 ns/op</td>
    <td>19877708 ns/op</td>
    <td>32% faster</td>
  </tr>
  <tr>
    <td rowspan="2" height="110px"><p align="center"><img src="https://www.gstatic.com/webp/gallery3/5.png" height="100px"></p></td>
    <td>file size</td>
    <td>139 kb</td>
    <td>119 kb</td>
    <td>14% smaller</td>
  </tr>
  <tr>
    <td>encoding time</td>
    <td>63423995 ns/op</td>
    <td>27813126 ns/op</td>
    <td>56% faster</td>
  </tr>
</table>
<p align="center">
<sub>image source: https://developers.google.com/speed/webp/gallery2</sub>
</p>


## Installation

To install the nativewebp package, use the following command:
```Bash
go get github.com/HugoSmits86/nativewebp
```
## Usage

Here’s a simple example of how to encode an image:
```Go
file, err := os.Create(name)
if err != nil {
  log.Fatalf("Error creating file %s: %v", name, err)
}
defer file.Close()

err = nativewebp.Encode(file, img, nil)
if err != nil {
  log.Fatalf("Error encoding image to WebP: %v", err)
}
```

Here’s a simple example of how to encode an animation:
```Go
file, err := os.Create(name)
if err != nil {
  log.Fatalf("Error creating file %s: %v", name, err)
}
defer file.Close()

ani := nativewebp.Animation{
  Images: []image.Image{
    frame1,
    frame2,
  },
  Durations: []uint {
    100,
    100,
  },
  Disposals: []uint {
    0,
    0,
  },
  LoopCount: 0,
  BackgroundColor: 0xffffffff,
}

err = nativewebp.EncodeAll(file, &ani, nil)
if err != nil {
  log.Fatalf("Error encoding WebP animation: %v", err)
}
```
