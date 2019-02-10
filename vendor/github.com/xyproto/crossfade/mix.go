package crossfade

import (
	"github.com/xyproto/imagelib"
	"image"
)

// Mixes two images with the given ratio. 0.5 is 50% of the first image and
// 50% of the second image. 0.1 is 10% of the first image and 90% of the
// second image.
func Files(inFilename1, inFilename2, outFilename string, ratio float64) error {
	img1, err := imagelib.Read(inFilename1)
	if err != nil {
		return err
	}
	img2, err := imagelib.Read(inFilename2)
	if err != nil {
		return err
	}

	// Crossfade
	SMul = ratio * 2.0
	DMul = -ratio*2.0 + 2.0
	outImage := BlendNewImage(img1, img2, Mix)

	err = imagelib.Write(outFilename, outImage)
	if err != nil {
		return err
	}
	return nil
}

// Experimental mix of two images, with a given ratio and contrast.
// Try ratio: 0.5 and contrast: 4.0.
func DemosceneMix(img1, img2 image.Image, ratio, contrast float64) image.Image {
	// Crossfade
	SMul = ratio * contrast
	DMul = 1.0 - ratio*contrast
	return BlendNewImage(img1, img2, OverlayMix)
}
