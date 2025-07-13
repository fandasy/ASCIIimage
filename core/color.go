package core

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
)

// Color represents color configuration for ASCII art rendering
//   - It ensures proper contrast between text (ascii char) and background
//   - When OriginalFace is true, it preserves the original pixel colors in output
type Color struct {
	// Face is the foreground/text color
	//  Ignored when OriginalFace is true.
	Face color.Color

	// Background is the canvas/background color
	Background color.Color

	// TransparentBackground removes the background
	TransparentBackground bool

	// OriginalFace preserves the source image colors
	OriginalFace bool

	// _Type caches the color model type for optimization.
	// Specifies the minimum color type to generate an image.
	_Type colorType
}

var (
	grayBlack = color.Gray{Y: 0x00}
	grayWhite = color.Gray{Y: 0xFF}
)

// DefaultColor returns the standard color scheme:
//   - Black text on white background
//   - TransparentBackground disabled (false)
//   - OriginalFace colors disabled (false)
//   - Uses grayscale colors for optimization
func DefaultColor() Color {
	return Color{
		Face:                  grayBlack,
		Background:            grayWhite,
		TransparentBackground: false,
		OriginalFace:          false,
	}
}

// validate ensures proper color configuration,
// It handles several special cases:
//   - When OriginalFace is true, only validates background
//   - Enforces contrast between Face and Background
//   - Replaces nil colors with complements
//   - Prevents identical Face/Background
//   - Converts colors to optimal format (_Type)
func (c *Color) validate() {
	var (
		faceNeed       = !c.OriginalFace
		backgroundNeed = !c.TransparentBackground
	)

	switch {
	case faceNeed && backgroundNeed:
		if c.Face == grayBlack && c.Background == grayWhite ||
			c.Face == grayWhite && c.Background == grayBlack {
			c._Type = colorTypeGray
			return
		}

		var (
			faceIsNil, _       = colorIsNilPtr(c.Face)
			backGroundIsNil, _ = colorIsNilPtr(c.Background)
		)

		switch {
		case faceIsNil && backGroundIsNil:
			c.Face = grayBlack
			c.Background = grayWhite
			c._Type = colorTypeGray
			return

		case faceIsNil:
			c.Face = complementaryColor(c.Background)

		case backGroundIsNil:
			c.Background = complementaryColor(c.Face)

		default:
			if c.Face == c.Background {
				c.Face = complementaryColor(c.Background)
			}
		}

		cType := getColorsType(c.Face, c.Background)
		c._Type = cType

		switch cType {
		case colorTypeGray:
			c.Face = colorToGray(c.Face)
			c.Background = colorToGray(c.Background)
		case colorTypeGray16:
			c.Face = colorToGray16(c.Face)
			c.Background = colorToGray16(c.Background)
		default:
			// ...
		}

	case faceNeed:
		faceIsNil, _ := colorIsNilPtr(c.Face)
		if faceIsNil {
			c.Face = grayBlack
			c._Type = colorTypeGray
		} else {
			c._Type = getColorType(c.Face)
		}

	case backgroundNeed:
		backGroundIsNil, _ := colorIsNilPtr(c.Background)
		if backGroundIsNil {
			c.Background = grayWhite
			c._Type = colorTypeGray
		} else {
			c._Type = getColorType(c.Face)
		}
	}
}

func (c *Color) isGray() bool {
	return c._Type == colorTypeGray || c._Type == colorTypeGray16
}

func (c *Color) createDrawImage(w, h int) draw.Image {
	switch {
	case c.TransparentBackground && c.OriginalFace:
		return image.NewRGBA(image.Rect(0, 0, w, h))
	case c.TransparentBackground:
		return image.NewRGBA(image.Rect(0, 0, w, h))
	case c.OriginalFace:
		return image.NewRGBA(image.Rect(0, 0, w, h))
	}

	switch c._Type {
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

func colorIsNilPtr(c color.Color) (bool, bool) {
	isNil := c == nil
	if isNil {
		return isNil, true
	}

	v := reflect.ValueOf(c)
	isPtr := v.Kind() == reflect.Ptr
	isNil = isPtr && v.IsNil()

	return isNil, isPtr
}

// complementaryColor generates an opposite color for maximum contrast.
func complementaryColor(c color.Color) color.Color {
	if v, ok := c.(color.Gray); ok {
		switch v {
		case grayWhite:
			return grayBlack
		case grayBlack:
			return grayWhite
		}
	}

	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(255 - uint(r>>8)),
		G: uint8(255 - uint(g>>8)),
		B: uint8(255 - uint(b>>8)),
		A: uint8(a),
	}
}

func colorToGray(c color.Color) color.Gray {
	switch v := c.(type) {
	case color.Gray:
		return v
	case *color.Gray:
		return *v
	}

	r, g, b, _ := c.RGBA()

	// (Rec. 709)
	y := (19595*r + 38470*g + 7471*b + 1<<15) >> 16

	return color.Gray{Y: uint8(y >> 8)}
}

func colorToGray16(c color.Color) color.Gray16 {
	switch v := c.(type) {
	case color.Gray16:
		return v
	case *color.Gray16:
		return *v
	}

	r, g, b, _ := c.RGBA()

	// (Rec. 709)
	y := (19595*r + 38470*g + 7471*b + 1<<15) >> 16

	return color.Gray16{Y: uint16(y)}
}
