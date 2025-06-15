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
}

// DefaultColor returns the standard color scheme:
//   - Black text on white background
func DefaultColor() Color {
	return Color{
		Face:       color.Black,
		Background: color.White,
	}
}

// validate ensures the color combination meets contrast requirements.
// It handles several special cases:
//   - Preserves standard black/white or white/black combinations
//   - Replaces nil colors with complementary colors
//   - Ensures foreground and background aren't identical
//   - Automatically generates complementary colors when needed
func (c *Color) validate() {
	if c.Face == color.Black && c.Background == color.White ||
		c.Face == color.White && c.Background == color.Black {
		return
	}

	var (
		faceIsNil       = colorIsNil(c.Face)
		backGroundIsNil = colorIsNil(c.Background)
	)

	switch {
	case faceIsNil && backGroundIsNil:
		c.Face = color.Black
		c.Background = color.White
		return
	case faceIsNil:
		c.Face = complementaryColor(c.Background)
		return
	case backGroundIsNil:
		c.Background = complementaryColor(c.Face)
		return
	}

	if c.Face == c.Background {
		c.Face = complementaryColor(c.Background)
	}
}

func colorIsNil(c color.Color) bool {
	return c == nil || func() bool {
		v := reflect.ValueOf(c)
		return v.Kind() == reflect.Ptr && v.IsNil()
	}()
}

// complementaryColor generates an opposite color for maximum contrast.
func complementaryColor(c color.Color) color.Color {
	if v, ok := c.(color.Gray16); ok {
		switch v {
		case color.White:
			return color.Black
		case color.Black:
			return color.White
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
