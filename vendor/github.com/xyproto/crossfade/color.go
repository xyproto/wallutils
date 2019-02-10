package crossfade

import (
	"image/color"
	"math"
)

type rgbaf64 struct {
	r, g, b, a float64
}

func (c rgbaf64) RGBA() (uint32, uint32, uint32, uint32) {
	r := float64ToUint16(c.r)
	g := float64ToUint16(c.g)
	b := float64ToUint16(c.b)
	a := float64ToUint16(c.a)
	return uint32(r), uint32(g), uint32(b), uint32(a)
}

type hslf64 struct {
	h, s, l float64
}

func (c hslf64) RGBA() (uint32, uint32, uint32, uint32) {
	return hsl2rgb(c.h, c.s, c.l).RGBA()
}

func color2rgbaf64(c color.Color) rgbaf64 {
	r, g, b, a := c.RGBA()
	return rgbaf64{float64(r), float64(g), float64(b), float64(a)}
}

func rgb2hsl(c color.Color) hslf64 {
	var h, s, l float64
	col := color2rgbaf64(c)
	r, g, b := col.r/max, col.g/max, col.b/max
	cmax := math.Max(math.Max(r, g), b)
	cmin := math.Min(math.Min(r, g), b)
	l = (cmax + cmin) / 2.0
	if cmax == cmin {
		// Achromatic.
		h, s = 0.0, 0.0
	} else {
		// Chromatic.
		delta := cmax - cmin
		if l > 0.5 {
			s = delta / (2.0 - cmax - cmin)
		} else {
			s = delta / (cmax + cmin)
		}
		switch cmax {
		case r:
			h = (g - b) / delta
			if g < b {
				h += 6.0
			}
		case g:
			h = (b-r)/delta + 2.0
		case b:
			h = (r-g)/delta + 4.0
		}
		h /= 6.0
	}
	return hslf64{h, s, l}
}

func hsl2rgb(h, s, l float64) color.Color {
	var r, g, b float64
	if s == 0 {
		r, g, b = l, l, l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - s*l
		}
		p := 2*l - q
		r = hue2rgb(p, q, h+1.0/3)
		g = hue2rgb(p, q, h)
		b = hue2rgb(p, q, h-1.0/3)
	}
	return rgbaf64{r*max + 0.5, g*max + 0.5, b*max + 0.5, max}
}

func hue2rgb(p, q, t float64) float64 {
	if t < 0.0 {
		t += 1.0
	}
	if t > 1.0 {
		t -= 1.0
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6.0*t
	}
	if t < 0.5 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6.0
	}
	return p
}
