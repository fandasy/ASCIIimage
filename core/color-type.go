package core

import (
	"image/color"
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

func getColorType(c color.Color) colorType {
	gray, alpha, depth := getColorAttr(c)

	switch {
	case gray && !alpha && depth == 8:
		return colorTypeGray
	case gray && !alpha && depth == 16:
		return colorTypeGray16
	case gray && alpha && depth == 8:
		return colorTypeNRGBA
	case gray && alpha && depth == 16:
		return colorTypeNRGBA64
	case !gray && !alpha && depth == 8:
		return colorTypeRGBA
	case !gray && !alpha && depth == 16:
		return colorTypeRGBA64
	case !gray && alpha && depth == 8:
		return colorTypeNRGBAc
	default: // !allGray && hasAlpha && bitDepth == 16
		return colorTypeNRGBA64c
	}
}

func getColorsType(c1, c2 color.Color) colorType {
	gray1, alpha1, depth1 := getColorAttr(c1)
	gray2, alpha2, depth2 := getColorAttr(c2)

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

func getColorAttr(c color.Color) (isGray bool, hasAlpha bool, bitDepth uint) {
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
