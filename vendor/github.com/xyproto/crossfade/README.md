# crossfade

* Go module for crossfading between two image files.
* Started out as a fork of [phrozen/blend](https://github.com/phrozen/blend).
* Includes a `blend` command line utility, for mixing two images 50%/50%.

## Example use of the Go package

    crossfade.Files("a.png", "b.png", "out.png", 0.5)

The last argument is a float that regulates the transition from one image to the other, where 0.0 is only `"a.png"`, while 1.0 is only `"b.png"`.

## Screenshots

0% lemur 100% mountain

![0% lemur](img/lagginhorn.jpg)

20% lemur 80% mountain

![20% lemur](img/out80.png)

50% lemur 50% mountain

![50% lemur](img/out50.png)

80% lemur 20% mountain

![80% lemur](img/out20.png)

100% lemur 0% mountain

![lemur](img/lemur.jpg)

The images are from wikipedia: <a href="https://en.wikipedia.org/wiki/File:Eulemur_mongoz_(male_-_face).jpg">lemur</a> | [mountain](https://nn.wikipedia.org/wiki/Fil:Lagginhorn_west_face.jpg)

## General info

* License: MIT
* Version: 2.0.0
