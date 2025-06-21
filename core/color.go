package core

import (
	"image/color"
	"reflect"
)

// Color represents a pair of foreground and background colors for ASCII art rendering.
// It ensures proper contrast between text and background.
type Color struct {
	// Face is the foreground/text color
	Face color.Color

	// Background is the canvas/background color
	Background color.Color

	// _Type defines what type of color.Color will be used when generating img
	_Type colorType
}

var (
	grayBlack = color.Gray{Y: 0x00}
	grayWhite = color.Gray{Y: 0xFF}
)

// DefaultColor returns the standard color scheme:
//   - Black text on white background
func DefaultColor() Color {
	return Color{
		Face:       grayBlack,
		Background: grayWhite,
	}
}

// validate ensures the color combination meets contrast requirements.
// It handles several special cases:
//   - Preserves standard black/white or white/black combinations
//   - Replaces nil colors with complementary colors
//   - Ensures foreground and background aren't identical
//   - Automatically generates complementary colors when needed
func (c *Color) validate() {
	if c.Face == grayBlack && c.Background == grayWhite ||
		c.Face == grayWhite && c.Background == grayBlack {
		c._Type = colorTypeGray
		return
	}

	var (
		faceIsNil       = colorIsNil(c.Face)
		backGroundIsNil = colorIsNil(c.Background)
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

	iType := getColorType(c.Face, c.Background)
	c._Type = iType
}

func colorIsNil(c color.Color) bool {
	return c == nil || func() bool {
		v := reflect.ValueOf(c)
		return v.Kind() == reflect.Ptr && v.IsNil()
	}()
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
