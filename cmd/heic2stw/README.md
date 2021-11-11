# heic2stw

Extract timing metadata from .heic (macOS dynamic wallpaper / HEIF) files.

Can be used in combination with an ImageMagick command like this:

    convert image.heic image_%02d.jpg

And then extract a simple timed wallpaper file like this:

    heic2stw image.heic > image.stw

Then copy `image.stw` and `image_*.jpg` to `/usr/share/backgrounds/image/`.

The included `heic-install` script does all of the above.
