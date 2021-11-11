# heif2stw

Extract timing metadata from .heic (macOS dynamic wallpaper / HEIF) files.

Can be used in combination with an ImageMagick command like this:

    convert example.heic example_%02d.jpg

And then extract a simple timed wallpaper file like this:

    heic2stw example.heic > example.stw

Then copy `example.stw` and `example_*.jpg` to `/usr/share/backgrounds/example/`.

For all of the above, use another word than `example` if your image has a different name.
