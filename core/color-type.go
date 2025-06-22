package core

import (
	"image"
	"image/color"
	"image/draw"
)

// colorType defines what type of color.Color will be used when generating img
type colorType int8

const (
	_                colorType = iota
	colorTypeGray              // image.Gray
	colorTypeGray16            // image.Gray16
	colorTypeNRGBA             // image.NRGBA (8-bit grayscale + alpha)
	colorTypeNRGBA64           // image.NRGBA64 (16-bit grayscale + alpha)
	colorTypeRGBA              // image.RGBA (8-bit color)
	colorTypeRGBA64            // image.RGBA64 (16-bit color)

	// The prefix c stands for color (not grayscale) (!allGray)
	colorTypeNRGBAc   // image.NRGBA (8-bit color + alpha)
	colorTypeNRGBA64c // image.NRGBA64 (16-bit color + alpha)
)

func getColorType(c1, c2 color.Color) colorType {
	gray1, alpha1, depth1 := getSingleColorType(c1)
	gray2, alpha2, depth2 := getSingleColorType(c2)

	allGray := gray1 && gray2
	hasAlpha := alpha1 || alpha2
	bitDepth := max(depth1, depth2)

	switch {
	case allGray && !hasAlpha && bitDepth == 8:
		return colorTypeGray
	case allGray && !hasAlpha && bitDepth == 16:
		return colorTypeGray16
	case allGray && hasAlpha && bitDepth == 8:
		return colorTypeNRGBA
	case allGray && hasAlpha && bitDepth == 16:
		return colorTypeNRGBA64
	case !allGray && !hasAlpha && bitDepth == 8:
		return colorTypeRGBA
	case !allGray && !hasAlpha && bitDepth == 16:
		return colorTypeRGBA64
	case !allGray && hasAlpha && bitDepth == 8:
		return colorTypeNRGBAc
	default: // !allGray && hasAlpha && bitDepth == 16
		return colorTypeNRGBA64c
	}
}

func getSingleColorType(c color.Color) (isGray bool, hasAlpha bool, bitDepth uint) {
	switch v := c.(type) {
	case color.Gray, *color.Gray:
		return true, false, 8
	case color.Gray16, *color.Gray16:
		return true, false, 16
	case color.RGBA:
		return v.R == v.G && v.G == v.B, false, 8
	case *color.RGBA:
		return v.R == v.G && v.G == v.B, false, 8
	case color.RGBA64:
		return v.R == v.G && v.G == v.B, false, 16
	case *color.RGBA64:
		return v.R == v.G && v.G == v.B, false, 16
	case color.NRGBA:
		return v.R == v.G && v.G == v.B, v.A != 255, 8
	case *color.NRGBA:
		return v.R == v.G && v.G == v.B, v.A != 255, 8
	case color.NRGBA64:
		return v.R == v.G && v.G == v.B, v.A != 65535, 16
	case *color.NRGBA64:
		return v.R == v.G && v.G == v.B, v.A != 65535, 16
	default:
		r, g, b, a := c.RGBA()
		isGray = r == g && g == b
		hasAlpha = a != 65535
		bitDepth = 16
		return
	}
}

func (ct colorType) createDrawImage(w, h int) draw.Image {
	switch ct {
	case colorTypeGray:
		return image.NewGray(image.Rect(0, 0, w, h))
	case colorTypeGray16:
		return image.NewGray16(image.Rect(0, 0, w, h))
	case colorTypeNRGBA:
		return image.NewNRGBA(image.Rect(0, 0, w, h))
	case colorTypeNRGBA64:
		return image.NewNRGBA64(image.Rect(0, 0, w, h))
	case colorTypeRGBA:
		return image.NewRGBA(image.Rect(0, 0, w, h))
	case colorTypeRGBA64:
		return image.NewRGBA64(image.Rect(0, 0, w, h))
	case colorTypeNRGBAc:
		return image.NewNRGBA(image.Rect(0, 0, w, h))
	case colorTypeNRGBA64c:
		return image.NewNRGBA64(image.Rect(0, 0, w, h))
	default:
		return image.NewNRGBA64(image.Rect(0, 0, w, h))
	}
}
