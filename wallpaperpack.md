# Ideas and thoughts

## Wallpaper Packs

### Wallpaper pack type 1, "simple", will set the picture that is there:

    wallpaper.zip
    |- filename.jpg

### Wallpaper pack type 2, "random", will change wallpaper every hour:

    wallpaper.zip
    |- filename1.jpg
    |- filename2.jpg
    |- filename3.jpg

### Wallpaper pack type 3, "hours", will change wallpaper based on the timestamp:

    wallpaper.zip
    |- 08:00-filename.jpg
    |- 10:00-filename.jpg
    |- 23:59-filename.jpg

### Wallpaper pack type 3, "hours", another example, using a date, the year is ignored:

    wallpaper.zip
    |- 01-04-2018-filename.jpg
    |- 12:30-filename.jpg

### Wallpaper pack type 4, "config", uses an expressive config file:

    wallpaper.zip
    |- conf.toml
    |- filename.jpg

    conf.toml:

    [event]
    at=08:00
    image = "filename.jpg"
    image.hue = "blue"
    once=false

    [event]
    at=default
    image = "filename.jpg"

---
    at=default means that no other event kicked in

### Wallpaper pack type 5, "plugin", uses Go `.so` files compiled with `go build -buildmode=plugin`.

    wallpaper.zip
    |- main.so
    |- filename1.jpg
    |- filename2.jpg

This needs a clearly defined interface fow interacting with the plugin, where the local time is sent to the plugin and pixels are sent back.
The plugin might want to draw something every minute. Live backgrounds?

This plugin support would make this module a "wallpaper engine".

### Wallpaper pack type 6, "script", uses an embedded scripting language, like Lua or Python

    wallpaper.zip
    |- main.lua
    |- filename.jpg

    main.lua:

    if hour == 8 then
        image = load("filename.jpg")
        image.hue = "blue"
        wallpaper.set(image)
    else
        image = load("filename.jpg")
        wallpaper.set(image)
    end

## Wallpaper pack type 7, "shared library", uses a .so file written in C or Rust

Good examples for C and Rust would be important.

## Global settings

* The user should be able to configure what type of wallpaper packs are okay to load.
* The "plugin" and "script" type should probably be disallowed by default?
* Or there should be some sort of voting or review system in place?
* Or it should be expected that wallpaper plugins installed into `/usr/share/wallpaper` are to be trusted.
* The user should be able to disable or disable any special holidays.

# Other considerations

* `/usr/lib/wallpaper` could be used to WM plugins, while `/usr/share/wallpaper/` could be used for wallpaper-packs.
* There could be a webpage for wallpaper packs, with preview images and preview animations (where 24-hours is sped up, and images for special dates, like halloween and christmas are displayed by themselves).
* The `https://github.com/xyproto/calendar` package could perhaps be used for identifying special dates, by string.

# Other animated desktop file formats

* edj, Enlightenment Desktop wallpaper pack, with animations

# Configuration

* A well chosen configuration file format, in combination with a small GUI client for modifying the settings, might be a good idea.
