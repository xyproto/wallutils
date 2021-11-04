# Simple Timed Wallpaper Format Spec

2019-03-12

A text format for specifying images and image transitions that make up timed wallpapers.

## Version 1.0.0

Simple timed wallpapers are UTF-8 encoded text files.

Every line may either start with `@`, for timing information, or with a field name followed by a colon `:` and a value.

### Key/value fields

The recognized fields are:

* `stw` (required), for specifying the version of the Simple Timed Wallpaper Format, for example `1.0`.
* `name` (optional), for giving the timed wallpaper a name.
* `format` (optional), for specifying a format string that may contain a `%s` marker. The format string will be used in the timing information.

After the fields, timing information may be specified. There are two types of timing information: static images or image transitions.

### Static images

Specifying a static image at a certain time, may look like this:

    @08:00: /usr/share/wallpapers/morning.jpg

This will change the wallpaper to `/usr/share/wallpapers/morning.jpg` when the event triggers at `08:00`.

Format description:

* The line must start with `@` followed by two digits which is the hour number.
* Then comes a colon `:` and two digits which is the minute number.
* Then comes a colon `:`, an optional whitespace, and a filename.
* The filename should not be quoted, and spaces in the filename are allowed, without any escaping.

Alternatively, a format string may be used. That would make the above example look like this:

    format: /usr/share/wallpapers/%s.jpg
    @08:00: morning

The `%s` marker will be replaced with the word `morning` when interpreting the filename for the `@08:00` event.

### Image transitions

Specifying an image transition, may look like this:

    @10:00-12:00: /usr/share/wallpapers/morning.jpg .. /usr/share/wallpapers/day.jpg | overlay

This will change the wallpaper to `/usr/share/wallpapers/morning.jpg` at `10:00`, then cross fade it to `/usr/share/wallpapers/day.jpg` in the 2 hours from `10:00` to `12:00` and the transition type will be `overlay`.

`overlay` is the default transition type and may be omitted. Implementing a cross fade between the first and second image is acceptable.

It is up to the implementation how often the wallpaper should be updated in the transition period from `10:00` to `12:00`. The recommendation is 10 times, regardless of the length of the time interval.

Format description:

* The line must start with `@` followed by two digits which is the hour number.
* Then comes a colon `:` and two digits which is the minute number.
* Then comes an optional whitespace, a dash `-` and another optional whitespace.
* Then comes two digits which is the hour number.
* Then comes a colon `:` and two digits which is the minute number.
* The first of the two timestamps is inclusive, while the second one is exclusive.
* Then comes a colon `:`, an optional whitespace, and an image filename that will be transitioned from.
* Then comes an optional whitespace, two dots `..` and another optional whitespace.
* Then comes an image filename that will be transitioned to.
* The filenames should not be quoted, and spaces in the filename are allowed, without any escaping.
* After the filenames, an optional space, a pipe `|`, an optional space and a transition type may be specified. This is optional.
* The only supported transition type for version 1.0 of the Simple Timed Wallpaper Format is `overlay`, which is also the default transition type.

Alternatively, a format string may be used. That would make the above example look like this:

    format: /usr/share/wallpapers/%s.jpg
    @10:00-12:00: morning .. day

## Real world examples

Two examples of GNOME Timed Wallpaper XML files converted to the Simple Timed Wallpaper format follows.

### mojave-timed

**mojave-timed.xml**

```xml
<background>
  <starttime>
    <year>2000</year>
    <month>01</month>
    <day>01</day>
    <hour>01</hour>
    <minute>00</minute>
    <second>00</second>
  </starttime>

  <transition type="overlay">
    <duration>14400.0</duration>
    <from>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-0100.jpg</from>
    <to>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-0500.jpg</to>
  </transition>
  <transition type="overlay">
    <duration>3600.0</duration>
    <from>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-0500.jpg</from>
    <to>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-0600.jpg</to>
  </transition>
  <transition type="overlay">
    <duration>3600.0</duration>
    <from>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-0600.jpg</from>
    <to>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-0700.jpg</to>
  </transition>
  <transition type="overlay">
    <duration>3600.0</duration>
    <from>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-0700.jpg</from>
    <to>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-0800.jpg</to>
  </transition>
  <transition type="overlay">
    <duration>3600.0</duration>
    <from>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-0800.jpg</from>
    <to>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-0900.jpg</to>
  </transition>
  <transition type="overlay">
    <duration>3600.0</duration>
    <from>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-0900.jpg</from>
    <to>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-1000.jpg</to>
  </transition>
  <transition type="overlay">
    <duration>3600.0</duration>
    <from>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-1000.jpg</from>
    <to>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-1100.jpg</to>
  </transition>
  <transition type="overlay">
    <duration>3600.0</duration>
    <from>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-1100.jpg</from>
    <to>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-1200.jpg</to>
  </transition>
  <transition type="overlay">
    <duration>4800.0</duration>
    <from>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-1200.jpg</from>
    <to>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-1320.jpg</to>
  </transition>
  <transition type="overlay">
    <duration>4800.0</duration>
    <from>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-1320.jpg</from>
    <to>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-1440.jpg</to>
  </transition>
  <transition type="overlay">
    <duration>4800.0</duration>
    <from>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-1440.jpg</from>
    <to>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-1600.jpg</to>
  </transition>
  <transition type="overlay">
    <duration>4800.0</duration>
    <from>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-1600.jpg</from>
    <to>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-1720.jpg</to>
  </transition>
  <transition type="overlay">
    <duration>4800.0</duration>
    <from>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-1720.jpg</from>
    <to>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-1840.jpg</to>
  </transition>
  <transition type="overlay">
    <duration>4800.0</duration>
    <from>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-1840.jpg</from>
    <to>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-2000.jpg</to>
  </transition>
  <transition type="overlay">
    <duration>3600.0</duration>
    <from>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-2000.jpg</from>
    <to>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-2100.jpg</to>
  </transition>
  <transition type="overlay">
    <duration>14400.0</duration>
    <from>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-2100.jpg</from>
    <to>/usr/share/backgrounds/gnome/mojave/mojave_dynamic-0100.jpg</to>
  </transition>
</background>
```

**mojave-timed.stw**

```yml
stw: 1.0
name: mojave-timed
format: /usr/share/backgrounds/gnome/mojave/mojave_dynamic-%s0.jpg
@01:00-05:00: 010 .. 050
@05:00-06:00: 050 .. 060
@06:00-07:00: 060 .. 070
@07:00-08:00: 070 .. 080
@08:00-09:00: 080 .. 090
@09:00-10:00: 090 .. 100
@10:00-11:00: 100 .. 110
@11:00-12:00: 110 .. 120
@12:00-13:20: 120 .. 132
@13:20-14:40: 132 .. 144
@14:40-16:00: 144 .. 160
@16:00-17:20: 160 .. 172
@17:20-18:40: 172 .. 184
@18:40-20:00: 184 .. 200
@20:00-21:00: 200 .. 210
@21:00-01:00: 210 .. 010
```

### adwaita-timed

**adwaita-timed.xml**

```xml
<background>
  <starttime>
    <year>2011</year>
    <month>11</month>
    <day>24</day>
    <hour>7</hour>
    <minute>00</minute>
    <second>00</second>
  </starttime>

<!-- This animation will start at 7 AM. -->

<!-- We start with sunrise at 7 AM. It will remain up for 1 hour. -->
<static>
<duration>3600.0</duration>
<file>/usr/share/backgrounds/gnome/adwaita-morning.jpg</file>
</static>

<!-- Sunrise starts to transition to day at 8 AM. The transition lasts for 5 hours, ending at 1 PM. -->
<transition type="overlay">
<duration>18000.0</duration>
<from>/usr/share/backgrounds/gnome/adwaita-morning.jpg</from>
<to>/usr/share/backgrounds/gnome/adwaita-day.jpg</to>
</transition>

<!-- It's 1 PM, we're showing the day image in full force now, for 5 hours ending at 6 PM. -->
<static>
<duration>18000.0</duration>
<file>/usr/share/backgrounds/gnome/adwaita-day.jpg</file>
</static>

<!-- It's 7 PM and it's going to start to get darker. This will transition for 6 hours up until midnight. -->
<transition type="overlay">
<duration>21600.0</duration>
<from>/usr/share/backgrounds/gnome/adwaita-day.jpg</from>
<to>/usr/share/backgrounds/gnome/adwaita-night.jpg</to>
</transition>

<!-- It's midnight. It'll stay dark for 5 hours up until 5 AM. -->
<static>
<duration>18000.0</duration>
<file>/usr/share/backgrounds/gnome/adwaita-night.jpg</file>
</static>

<!-- It's 5 AM. We'll start transitioning to sunrise for 2 hours up until 7 AM. -->
<transition type="overlay">
<duration>7200.0</duration>
<from>/usr/share/backgrounds/gnome/adwaita-night.jpg</from>
<to>/usr/share/backgrounds/gnome/adwaita-morning.jpg</to>
</transition>
</background>
```

**adwaita-timed.stw**

```yml
stw: 1.0
name: adwaita-timed
format: /usr/share/backgrounds/gnome/adwaita-%s.jpg
@07:00: morning
@08:00-13:00: morning .. day
@13:00: day
@18:00-00:00: day .. night
@00:00: night
@05:00-07:00: night .. morning
```

### Final remarks

The `xml2stw` utility can be used for converting GNOME timed XML files to the Simple Timed Wallpaper format. It was used for converting the examples above.

This is a draft. Pull requests are welcome: https://github.com/xyproto/simpletimed/pulls

### Author

Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
