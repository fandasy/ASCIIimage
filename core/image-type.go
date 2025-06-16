package core

import (
	"image"
	"image/color"
	"image/draw"
)

// imageType defines what type of color.Color will be used when generating img
type imageType int8

const (
	_                imageType = iota
	imageTypeGray              // image.Gray
	imageTypeGray16            // image.Gray16
	imageTypeNRGBA             // image.NRGBA (8-bit grayscale + alpha)
	imageTypeNRGBA64           // image.NRGBA64 (16-bit grayscale + alpha)
	imageTypeRGBA              // image.RGBA (8-bit color)
	imageTypeRGBA64            // image.RGBA64 (16-bit color)

	// The prefix c stands for color (not grayscale) (!allGray)
	imageTypeNRGBAc   // image.NRGBA (8-bit color + alpha)
	imageTypeNRGBA64c // image.NRGBA64 (16-bit color + alpha)
)

func getColorType(c1, c2 color.Color) imageType {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	hasAlpha := a1 != 0xFFFF || a2 != 0xFFFF

	isGray1 := r1 == g1 && g1 == b1
	isGray2 := r2 == g2 && g2 == b2
	allGray := isGray1 && isGray2

	bitDepth := uint(8)
	if r1 > 0xFF || g1 > 0xFF || b1 > 0xFF || a1 > 0xFF ||
		r2 > 0xFF || g2 > 0xFF || b2 > 0xFF || a2 > 0xFF {
		bitDepth = 16
	}

	switch {
	case allGray && !hasAlpha && bitDepth == 8:
		return imageTypeGray
	case allGray && !hasAlpha && bitDepth == 16:
		return imageTypeGray16
	case allGray && hasAlpha && bitDepth == 8:
		return imageTypeNRGBA
	case allGray && hasAlpha && bitDepth == 16:
		return imageTypeNRGBA64
	case !allGray && !hasAlpha && bitDepth == 8:
		return imageTypeRGBA
	case !allGray && !hasAlpha && bitDepth == 16:
		return imageTypeRGBA64
	case !allGray && hasAlpha && bitDepth == 8:
		return imageTypeNRGBAc
	default: // !allGray && hasAlpha && bitDepth == 16
		return imageTypeNRGBA64c
	}
}

func (it imageType) createDrawImage(w, h int) draw.Image {
	switch it {
	case imageTypeGray:
		return image.NewGray(image.Rect(0, 0, w, h))
	case imageTypeGray16:
		return image.NewGray16(image.Rect(0, 0, w, h))
	case imageTypeNRGBA:
		return image.NewNRGBA(image.Rect(0, 0, w, h))
	case imageTypeNRGBA64:
		return image.NewNRGBA64(image.Rect(0, 0, w, h))
	case imageTypeRGBA:
		return image.NewRGBA(image.Rect(0, 0, w, h))
	case imageTypeRGBA64:
		return image.NewRGBA64(image.Rect(0, 0, w, h))
	case imageTypeNRGBAc:
		return image.NewNRGBA(image.Rect(0, 0, w, h))
	case imageTypeNRGBA64c:
		return image.NewNRGBA64(image.Rect(0, 0, w, h))
	default:
		return image.NewNRGBA64(image.Rect(0, 0, w, h))
	}
}
