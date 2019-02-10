// Package imagelib can deal with RGB, HSV, HSL and PNG images
package imagelib

import (
	"image"
	"image/color"
	"math"
)

// Pick out only the red colors from an image
func Red(m image.Image) image.Image {
	var (
		rect     = m.Bounds()
		c        color.Color
		cr       color.RGBA
		newImage = image.NewRGBA(image.Rect(0, 0, rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y))
	)
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			c = m.At(x, y)      // c is RGBAColor, which implements Color
			cr = c.(color.RGBA) // this is needed
			c = color.RGBA{cr.R, 0, 0, cr.A}
			newImage.Set(x-rect.Min.X, y-rect.Min.Y, c)
		}
	}
	return newImage
}

// Pick out only the green colors from an image
func Green(m image.Image) image.Image {
	var (
		rect     = m.Bounds()
		c        color.Color
		cr       color.RGBA
		newImage = image.NewRGBA(image.Rect(0, 0, rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y))
	)
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			c = m.At(x, y)      // c is RGBAColor, which implements Color
			cr = c.(color.RGBA) // this is needed
			c = color.RGBA{0, cr.G, 0, cr.A}
			newImage.Set(x-rect.Min.X, y-rect.Min.Y, c)
		}
	}
	return newImage
}

// Pick out only the blue colors from an image
func Blue(m image.Image) image.Image {
	var (
		rect     = m.Bounds()
		c        color.Color
		cr       color.RGBA
		newImage = image.NewRGBA(image.Rect(0, 0, rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y))
	)
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			c = m.At(x, y)      // c is RGBAColor, which implements Color
			cr = c.(color.RGBA) // this is needed
			c = color.RGBA{0, 0, cr.B, cr.A}
			newImage.Set(x-rect.Min.X, y-rect.Min.Y, c)
		}
	}
	return newImage
}

// Pick out only the colors close to the given color,
// within a given threshold
func CloseTo1(m image.Image, target color.RGBA, thresh uint8) image.Image {
	var (
		rect       = m.Bounds()
		c          color.Color
		cr         color.RGBA
		r, g, b, a uint8
		newImage   = image.NewRGBA(image.Rect(0, 0, rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y))
	)
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			c = m.At(x, y)      // c is RGBAColor, which implements Color
			cr = c.(color.RGBA) // this is needed
			r = 0
			g = 0
			b = 0
			a = cr.A
			if abs(int8(target.R)-int8(cr.R)) < thresh {
				r = target.R
			}
			if abs(int8(target.G)-int8(cr.G)) < thresh {
				g = target.G
			}
			if abs(int8(target.B)-int8(cr.B)) < thresh {
				b = target.B
			}
			c = color.RGBA{r, g, b, a}
			newImage.Set(x-rect.Min.X, y-rect.Min.Y, c)
		}
	}
	return newImage
}

// Pick out only the colors close to the given color,
// within a given threshold. Make it uniform.
// Zero alpha to unused pixels in returned image.
func CloseTo2(m image.Image, target color.RGBA, thresh uint8) image.Image {
	var (
		rect       = m.Bounds()
		c          color.Color
		cr         color.RGBA
		r, g, b, a uint8
		newImage   = image.NewRGBA(image.Rect(0, 0, rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y))
	)
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			c = m.At(x, y)      // c is RGBAColor, which implements Color
			cr = c.(color.RGBA) // this is needed
			r = 0
			g = 0
			b = 0
			a = 0
			if abs(int8(target.R)-int8(cr.R)) < thresh || abs(int8(target.G)-int8(cr.G)) < thresh || abs(int8(target.B)-int8(cr.B)) < thresh {
				r = target.R
				g = target.G
				b = target.B
				a = cr.A
			}
			c = color.RGBA{r, g, b, a}
			newImage.Set(x-rect.Min.X, y-rect.Min.Y, c)
		}
	}
	return newImage
}

// Take orig, add the nontransparent colors from addimage, as addascolor
func AddToAs(orig image.Image, addimage image.Image, addcolor color.RGBA) image.Image {
	var (
		rect       = addimage.Bounds()
		c          color.Color
		cr, or     color.RGBA
		r, g, b, a uint8
		newImage   = image.NewRGBA(image.Rect(0, 0, rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y))
	)
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			cr = addimage.At(x, y).(color.RGBA)
			or = orig.At(x, y).(color.RGBA)
			r = or.R
			g = or.G
			b = or.B
			a = or.A
			if cr.A > 0 {
				r = addcolor.R
				g = addcolor.G
				b = addcolor.B
				a = addcolor.A
			}
			c = color.RGBA{r, g, b, a}
			newImage.Set(x-rect.Min.X, y-rect.Min.Y, c)
		}
	}
	return newImage
}

// Convert RGB to hue
func Hue(cr color.RGBA) float64 {
	r := float64(cr.R) / 255.0
	g := float64(cr.G) / 255.0
	b := float64(cr.B) / 255.0
	var h float64
	RGBmax := r
	if g > RGBmax {
		RGBmax = g
	}
	if b > RGBmax {
		RGBmax = b
	}
	if RGBmax == r {
		h = 60 * (g - b)
		if h < 0 {
			h += 360
		}
	} else if RGBmax == g {
		h = 120 + 60*(b-r)
	} else /* RGBmax == rgb.b */ {
		h = 240 + 60*(r-g)
	}
	return h
}

// Convert RGB to HSV
func HSV(cr color.RGBA) (uint8, uint8, uint8) {
	var hue, sat, val uint8
	RGBmin := min(cr.R, cr.G, cr.B)
	RGBmax := max(cr.R, cr.G, cr.B)

	val = RGBmax
	if val == 0 {
		hue = 0
		sat = 0
		return hue, sat, val
	}

	sat = 255 * (RGBmax - RGBmin) / val
	if sat == 0 {
		hue = 0
		return hue, sat, val
	}

	span := (RGBmax - RGBmin)
	if RGBmax == cr.R {
		hue = 43 * (cr.G - cr.B) / span
	} else if RGBmax == cr.G {
		hue = 85 + 43*(cr.B-cr.R)/span
	} else { /* RGBmax == cr.B */
		hue = 171 + 43*(cr.R-cr.G)/span
	}

	return hue, sat, val
}

// Separate an image into three colors, with a given threshold
func Separate(inImage image.Image, color1, color2, color3 color.RGBA, thresh uint8, t float64) []image.Image {
	var (
		rect       = inImage.Bounds()
		cr         color.RGBA
		r, g, b, a uint8
		h, s       float64
		images     = make([]image.Image, 3) // 3 is the number of images
		newImage1  = image.NewRGBA(image.Rect(0, 0, rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y))
		newImage2  = image.NewRGBA(image.Rect(0, 0, rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y))
		newImage3  = image.NewRGBA(image.Rect(0, 0, rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y))
	)
	hue1, _, s1 := HLS(float64(color1.R)/255.0, float64(color1.G)/255.0, float64(color1.B)/255.0)
	hue2, _, s2 := HLS(float64(color2.R)/255.0, float64(color2.G)/255.0, float64(color2.B)/255.0)
	hue3, _, s3 := HLS(float64(color3.R)/255.0, float64(color3.G)/255.0, float64(color3.B)/255.0)
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			// get the rgba color
			// cr = inImage.At(x, y).(image.RGBAColor)
			cr = color.RGBAModel.Convert(inImage.At(x, y)).(color.RGBA)
			r = 0
			g = 0
			b = 0
			a = 255
			h, _, s = HLS(float64(cr.R)/255.0, float64(cr.G)/255.0, float64(cr.B)/255.0)
			// Find the closest color of the three, measured in hue and saturation
			if ((fabs(h-hue1) < fabs(h-hue2)) && (fabs(h-hue1) < fabs(h-hue3))) ||
				((fabs(s-s1) < fabs(s-s2)) && (fabs(s-s1) < fabs(s-s3))) {
				// Only add if the color is close enough
				if abs(int8(color1.R)-int8(cr.R)) < thresh || abs(int8(color1.G)-int8(cr.G)) < thresh || abs(int8(color1.B)-int8(cr.B)) < thresh {
					r = color1.R
					g = color1.G
					b = color1.B
					newImage1.Set(x-rect.Min.X, y-rect.Min.Y, color.RGBA{r, g, b, a})
				}
			} else if ((fabs(h-hue2) < fabs(h-hue1)) && (fabs(h-hue2) < fabs(h-hue3))) ||
				((fabs(s-s2) < fabs(s-s1)) && (fabs(s-s2) < fabs(s-s3))) {
				// Only add if the color is close enough
				if abs(int8(color2.R)-int8(cr.R)) < thresh || abs(int8(color2.G)-int8(cr.G)) < thresh || abs(int8(color2.B)-int8(cr.B)) < thresh {
					r = color2.R
					g = color2.G
					b = color2.B
					newImage2.Set(x-rect.Min.X, y-rect.Min.Y, color.RGBA{r, g, b, a})
				}
			} else if ((fabs(h-hue3) < fabs(h-hue1)) && (fabs(h-hue3) < fabs(h-hue2))) ||
				((fabs(s-s3) < fabs(s-s1)) && (fabs(s-s3) < fabs(s-s2))) {
				if abs(int8(color3.R)-int8(cr.R)) < thresh || abs(int8(color3.G)-int8(cr.G)) < thresh || abs(int8(color3.B)-int8(cr.B)) < thresh {
					r = color3.R
					g = color3.G
					b = color3.B
					newImage3.Set(x-rect.Min.X, y-rect.Min.Y, color.RGBA{r, g, b, a})
				}
			}
		}
	}
	images[0] = newImage1
	images[1] = newImage2
	images[2] = newImage3
	return images
}

// Convert RGB to HLS
func HLS(r, g, b float64) (float64, float64, float64) {
	// Ported from Python colorsys
	var h, l, s float64
	maxc := fmax(r, g, b)
	minc := fmin(r, g, b)
	l = (minc + maxc) / 2.0
	if minc == maxc {
		return 0.0, l, 0.0
	}
	span := (maxc - minc)
	if l <= 0.5 {
		s = span / (maxc + minc)
	} else {
		s = span / (2.0 - maxc - minc)
	}
	rc := (maxc - r) / span
	gc := (maxc - g) / span
	bc := (maxc - b) / span
	if r == maxc {
		h = bc - gc
	} else if g == maxc {
		h = 2.0 + rc - bc
	} else {
		h = 4.0 + gc - rc
	}
	h = math.Mod((h / 6.0), 1.0)
	return h, l, s
}

// Ported from Python colorsys
func _v(m1, m2, hue float64) float64 {
	oneSixth := 1.0 / 6.0
	twoThird := 2.0 / 3.0
	hue = math.Mod(hue, 1.0)
	if hue < oneSixth {
		return m1 + (m2-m1)*hue*6.0
	}
	if hue < 0.5 {
		return m2
	}
	if hue < twoThird {
		return m1 + (m2-m1)*(twoThird-hue)*6.0
	}
	return m1
}

// Convert a HLS color to RGB
func HLStoRGB(h, l, s float64) (float64, float64, float64) {
	// Ported from Python colorsys
	oneThird := 1.0 / 3.0
	if s == 0.0 {
		return l, l, l
	}
	var m2 float64
	if l <= 0.5 {
		m2 = l * (1.0 + s)
	} else {
		m2 = l + s - (l * s)
	}
	m1 := 2.0*l - m2
	return _v(m1, m2, h+oneThird), _v(m1, m2, h), _v(m1, m2, h-oneThird)
}

// Mix two RGB colors, a bit like how paint mixes
func PaintMix(c1, c2 color.RGBA) color.RGBA {
	// Thanks to Mark Ransom via stackoverflow

	// The less pi-presition, the greener the mix between blue and yellow
	// Using math.Pi gives a completely different result, for some reason
	//pi := math.Pi
	//pi := 3.141592653589793
	pi := 3.141592653
	//pi := 3.1415

	h1, l1, s1 := HLS(float64(c1.R)/255.0, float64(c1.G)/255.0, float64(c1.B)/255.0)
	h2, l2, s2 := HLS(float64(c2.R)/255.0, float64(c2.G)/255.0, float64(c2.B)/255.0)
	h := 0.0
	s := 0.5 * (s1 + s2)
	l := 0.5 * (l1 + l2)
	x := math.Cos(2.0*pi*h1) + math.Cos(2.0*pi*h2)
	y := math.Sin(2.0*pi*h1) + math.Sin(2.0*pi*h2)
	if (x != 0.0) || (y != 0.0) {
		h = math.Atan2(y, x) / (2.0 * pi)
	} else {
		s = 0.0
	}
	r, g, b := HLStoRGB(h, l, s)
	return color.RGBA{uint8(r * 255.0), uint8(g * 255.0), uint8(b * 255.0), 255}
}
