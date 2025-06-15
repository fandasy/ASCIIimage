package core

import "image/color"

// Options configure the ASCII art generation process
type Options struct {
	// PixelRatio defines how many original pixels map to one ASCII character
	// Format: X (width), Y (height) original pixels → 1 ASCII character
	PixelRatio PixelRatio

	// Chars defines the character set to use for brightness mapping
	Chars *Chars

	// Color specifies the foreground and background color scheme
	// If invalid or unset, defaults to black-on-white
	// Use DefaultColor() for standard scheme
	Color Color
}

// DefaultOptions returns the default conversion options:
//   - PixelRatio: 1x1 (one source pixel per ASCII character)
//   - Chars: Default character set ("@%#*+=:~-.  ")
//   - Color: Black text on white background
func DefaultOptions() *Options {
	return &Options{
		PixelRatio: DefaultPixelRatio(),
		Chars:      DefaultChars(),
		Color:      DefaultColor(),
	}
}

func (o *Options) WithPixelRatio(x, y int) *Options {
	o.PixelRatio = PixelRatio{X: x, Y: y}
	return o
}

func (o *Options) WithChars(c *Chars) *Options {
	o.Chars = c
	return o
}

func (o *Options) WithColor(color Color) *Options {
	o.Color = color
	return o
}

func (o *Options) WithFaceColor(c color.Color) *Options {
	o.Color.Face = c
	return o
}

func (o *Options) WithBackgroundColor(c color.Color) *Options {
	o.Color.Background = c
	return o
}

// validate ensures the options have valid values, setting defaults where needed
func (o *Options) validate() {
	if o.PixelRatio.X <= 0 {
		o.PixelRatio.X = defaultXPixelRatio
	}
	if o.PixelRatio.Y <= 0 {
		o.PixelRatio.Y = defaultYPixelRatio
	}

	if o.Chars == nil {
		o.Chars = DefaultChars()
	}

	o.Color.validate()
}
