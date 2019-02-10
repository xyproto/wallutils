// Package crossfade implements blending mode functions between images,
// and some utility functions for image processing.
//
// The fundamental part of the library is the type BlendFunc,
// the function is applied to each pixel where the top layer (src)
// overlaps the bottom layer (dst) of both given 'image' interfaces.
//
// This library provides many of the widely used blending functions
// to be used either as 'mode' parameter to the Blend() primary
// function, or to be used individually providing two 'color' interfaces.
// You can implement your own blending modes and pass them into the
// Blend() function.
//
// This is the list of the currently implemented blending modes:
//
// Add, Color, Color Burn, Color Dodge, Darken, Darker Color, Difference,
// Divide, Exclusion, Hard Light, Hard Mix, Hue, Lighten, Lighter Color,
// Linear Burn, Linear Dodge, Linear Light, Luminosity, Multiply, Overlay,
// Phoenix, Pin Light, Reflex, Saturation, Screen, Soft Light, Substract,
// Vivid Light.
package crossfade

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// Constants of max and mid values for uint16 for internal use.
// This can be changed to make the algorithms use uint8 instead,
// but they are kept this way to provide more acurate calculations
// and to support all of the color modes in the 'image' package.
const (
	max = 65535.0 // equals to 0xFFFF uint16 max range of color.Color
	mid = max / 2.0
)

var (
	Add          BlendFunc
	Color        BlendFunc
	ColorBurn    BlendFunc
	ColorDodge   BlendFunc
	Darken       BlendFunc
	DarkerColor  BlendFunc
	Difference   BlendFunc
	Divide       BlendFunc
	Exclusion    BlendFunc
	HardLight    BlendFunc
	HardMix      BlendFunc
	Hue          BlendFunc
	Lighten      BlendFunc
	LighterColor BlendFunc
	LinearBurn   BlendFunc
	LinearDodge  BlendFunc
	LinearLight  BlendFunc
	Luminosity   BlendFunc
	Multiply     BlendFunc
	Overlay      BlendFunc
	Mix          BlendFunc
	OverlayMix   BlendFunc
	Phoenix      BlendFunc
	PinLight     BlendFunc
	Reflex       BlendFunc
	Saturation   BlendFunc
	Screen       BlendFunc
	SoftLight    BlendFunc
	Substract    BlendFunc
	VividLight   BlendFunc
)

// A blend function or blend mode receives a destination color and
// a source color, then returns a transformation of them. Blend()
// function receives a BlendFunc and applies it to every pixel in
// the overlaping areas of two given images.
type BlendFunc func(dst, src color.Color) color.Color

// Blends src image (top layer) into dst image (bottom layer) using
// the BlendFunc provided by mode. BlendFunc is applied to each pixel
// where the src image overlaps the dst image and the result is stored
// in the original dst image, src image is unmutable.
func BlendImage(dst draw.Image, src image.Image, mode BlendFunc) {
	// Obtain the intersection of both images.
	inter := dst.Bounds().Intersect(src.Bounds())
	// Apply BlendFuc to each pixel in the intersection.
	for y := inter.Min.Y; y < inter.Max.Y; y++ {
		for x := inter.Min.X; x < inter.Max.X; x++ {
			dst.Set(x, y, mode(dst.At(x, y), src.At(x, y)))
		}
	}
}

// Blends src image (top layer) into dst image (bottom layer) using
// the BlendFunc provided by mode. BlendFunc is applied to each pixel
// where the src image overlaps the dst image and returns the resulting
// image without modifying src, or dst as they are both unmutable.
func BlendNewImage(dst, src image.Image, mode BlendFunc) image.Image {
	// Obtain the intersection of both images.
	inter := dst.Bounds().Intersect(src.Bounds())
	// Create a new RGBA or RGBA64 image to return the values.
	img := image.NewRGBA(dst.Bounds())
	// Iterate over dst image pixels.
	for y := dst.Bounds().Min.Y; y < dst.Bounds().Max.Y; y++ {
		for x := dst.Bounds().Min.X; x < dst.Bounds().Max.X; x++ {
			// If src is inside the intersection, we blend both
			// pixels using the provided BlendFunc (mode).
			if p := image.Pt(x, y); p.In(inter) {
				img.Set(x, y, mode(dst.At(x, y), src.At(x, y)))
			} else {
				// Else we copy dst pixel to the resulting image.
				img.Set(x, y, dst.At(x, y))
			}
		}
	}
	return img
}

func blendPerChannel(dst, src color.Color, bf func(float64, float64) float64) color.Color {
	d, s := color2rgbaf64(dst), color2rgbaf64(src)
	return rgbaf64{bf(d.r, s.r), bf(d.g, s.g), bf(d.b, s.b), d.a}
}

// Blending modes supported by Photoshop in order.
/*-------------------------------------------------------*/

// DARKEN
func darken(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, darken_per_ch)
}
func darken_per_ch(d, s float64) float64 {
	return math.Min(d, s)
}

// MULTIPLY
func multiply(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, multiply_per_ch)
}
func multiply_per_ch(d, s float64) float64 {
	return s * d / max
}

// COLOR BURN
func color_burn(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, color_burn_per_ch)
}
func color_burn_per_ch(d, s float64) float64 {
	if s == 0.0 {
		return s
	}
	return math.Max(0.0, max-((max-d)*max/s))
}

// LINEAR BURN
func linear_burn(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, linear_burn_per_ch)
}
func linear_burn_per_ch(d, s float64) float64 {
	if (s + d) < max {
		return 0.0
	}
	return s + d - max
}

// DARKER COLOR
func darker_color(dst, src color.Color) color.Color {
	s, d := color2rgbaf64(src), color2rgbaf64(dst)
	if s.r+s.g+s.b > d.r+d.g+d.b {
		return dst
	}
	return src
}

/*-------------------------------------------------------*/

// LIGHTEN
func lighten(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, lighten_per_ch)
}
func lighten_per_ch(d, s float64) float64 {
	return math.Max(d, s)
}

// SCREEN
func screen(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, screen_per_ch)
}
func screen_per_ch(d, s float64) float64 {
	return s + d - s*d/max
}

// COLOR DODGE
func color_dodge(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, color_dodge_per_ch)
}
func color_dodge_per_ch(d, s float64) float64 {
	if s == max {
		return s
	}
	return math.Min(max, (d * max / (max - s)))
}

// LINEAR DODGE
func linear_dodge(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, linear_dodge_per_ch)
}
func linear_dodge_per_ch(d, s float64) float64 {
	return math.Min(s+d, max)
}

// LIGHTER COLOR
func lighter_color(dst, src color.Color) color.Color {
	s, d := color2rgbaf64(src), color2rgbaf64(dst)
	if s.r+s.g+s.b > d.r+d.g+d.b {
		return src
	}
	return dst
}

/*-------------------------------------------------------*/

// OVERLAY
func overlay(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, overlay_per_ch)
}
func overlay_per_ch(d, s float64) float64 {
	if d < mid {
		return 2 * s * d / max
	}
	return max - 2*(max-s)*(max-d)/max
}

// Used by MIX and OVERLAY MIX to multiply with the source and destination
// color values
var (
	SMul = 1.0
	DMul = 1.0
)

// OVERLAY MIX
func overlay_mix(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, overlay_mix_per_ch)
}
func overlay_mix_per_ch(d, s float64) float64 {
	if d < mid {
		return 2 * (((s * SMul) + (d * DMul)) / 2.0) / max
	}
	return max - 2*(max-((s*SMul)/2.0))*(max-((d*DMul)/2.0))/max
}

// MIX
func mix(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, mix_per_ch)
}
func mix_per_ch(d, s float64) float64 {
	return ((s * SMul) + (d * DMul)) / 2.0
}

// SOFT LIGHT
func soft_light(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, soft_light_per_ch)
}
func soft_light_per_ch(d, s float64) float64 {
	return (d / max) * (d + (2*s/max)*(max-d))
}

// HARD LIGHT
func hard_light(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, hard_light_per_ch)
}
func hard_light_per_ch(d, s float64) float64 {
	if s > mid {
		return d + (max-d)*((s-mid)/mid)
	}
	return d * s / mid
}

// VIVID LIGHT (check)
func vivid_light(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, vivid_light_per_ch)
}
func vivid_light_per_ch(d, s float64) float64 {
	if s < mid {
		return color_burn_per_ch((2 * s), d)
	}
	return color_dodge_per_ch((2 * (s - mid)), d)
}

// LINEAR LIGHT
func linear_light(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, linear_light_per_ch)
}
func linear_light_per_ch(d, s float64) float64 {
	if s < mid {
		return linear_burn_per_ch((2 * s), d)
	}
	return linear_dodge_per_ch((2 * (s - mid)), d)
}

// PIN LIGHT
func pin_light(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, pin_light_per_ch)
}
func pin_light_per_ch(d, s float64) float64 {
	if s < mid {
		return darken_per_ch((2 * s), d)
	}
	return lighten_per_ch((2 * (s - mid)), d)
}

// HARD MIX (check)
func hard_mix(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, hard_mix_per_ch)
}
func hard_mix_per_ch(d, s float64) float64 {
	if vivid_light_per_ch(d, s) < mid {
		return 0.0
	}
	return max
}

/*-------------------------------------------------------*/

// DIFFERENCE
func difference(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, difference_per_ch)
}
func difference_per_ch(d, s float64) float64 {
	return math.Abs(s - d)
}

// EXCLUSION
func exclusion(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, exclusion_per_ch)
}
func exclusion_per_ch(d, s float64) float64 {
	return s + d - s*d/mid
}

// SUBSTRACT
func substract(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, substract_per_ch)
}
func substract_per_ch(d, s float64) float64 {
	if d-s < 0.0 {
		return 0.0
	}
	return d - s
}

// DIVIDE
func divide(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, divide_per_ch)
}
func divide_per_ch(d, s float64) float64 {
	return (d*max)/s + 1.0
}

// Blending modes that use HSL color model transformations.
/*-------------------------------------------------------*/

// HUE
func hue(dst, src color.Color) color.Color {
	s := rgb2hsl(src)
	if s.s == 0.0 {
		return dst
	}
	d := rgb2hsl(dst)
	return hsl2rgb(s.h, d.s, d.l)
}

// SATURATION
func saturation(dst, src color.Color) color.Color {
	s := rgb2hsl(src)
	d := rgb2hsl(dst)
	return hsl2rgb(d.h, s.s, d.l)
}

// COLOR "added _ to avoid namespace conflict with 'color' package"
func color_(dst, src color.Color) color.Color {
	s := rgb2hsl(src)
	d := rgb2hsl(dst)
	return hsl2rgb(s.h, s.s, d.l)
}

// LUMINOSITY
func luminosity(dst, src color.Color) color.Color {
	s := rgb2hsl(src)
	d := rgb2hsl(dst)
	return hsl2rgb(d.h, d.s, s.l)
}

// This blending modes are not implemented in Photoshop
// or GIMP at the moment, but produced their desired results.
/*-------------------------------------------------------*/

// ADD
func add(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, add_per_ch)
}
func add_per_ch(d, s float64) float64 {
	if s+d > max {
		return max
	}
	return s + d
}

// REFLEX (a.k.a GLOW)
func reflex(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, reflex_per_ch)
}
func reflex_per_ch(d, s float64) float64 {
	if s == max {
		return s
	}
	return math.Min(max, (d * d / (max - s)))
}

// PHOENIX
func phoenix(dst, src color.Color) color.Color {
	return blendPerChannel(dst, src, phoenix_per_ch)
}
func phoenix_per_ch(d, s float64) float64 {
	return math.Min(d, s) - math.Max(d, s) + max
}

// Init function maps the blending mode functions.
func init() {
	Darken = darken
	Multiply = multiply
	ColorBurn = color_burn
	LinearBurn = linear_burn
	DarkerColor = darker_color
	Lighten = lighten
	Screen = screen
	ColorDodge = color_dodge
	LinearDodge = linear_dodge
	LighterColor = lighter_color
	Overlay = overlay
	Mix = mix
	OverlayMix = overlay_mix
	SoftLight = soft_light
	HardLight = hard_light
	VividLight = vivid_light
	LinearLight = linear_light
	PinLight = pin_light
	HardMix = hard_mix
	Difference = difference
	Exclusion = exclusion
	Substract = substract
	Divide = divide
	Hue = hue
	Saturation = saturation
	Color = color_
	Luminosity = luminosity
	Add = add
	Reflex = reflex
	Phoenix = phoenix
}
