# Simpletimed

`simpletimed` can be used for parsing the STW file format and for running an event loop for setting the wallpaper, given a function with this signature:

```go
func(string) error
```

Where the given string is the image filename to be set.

## The Simple Timed Wallpaper Format

STW is a format for a configuration file that specifies in which time ranges wallpapers should change from one to another, and with which transition.

It's a similar to the GNOME timed wallpaper XML format, but much simpler and less verbose.
