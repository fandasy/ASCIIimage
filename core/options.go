package core

// Options configure the ASCII art generation process
type Options struct {
	// PixelRatio defines how many original pixels map to one ASCII character
	// Format: X (width), Y (height) original pixels → 1 ASCII character
	PixelRatio PixelRatio

	// Chars defines the character set to use for brightness mapping
	Chars *Chars
}

// DefaultOptions returns the default conversion options:
// - PixelRatio: 1x1 (one source pixel per ASCII character)
// - Chars: Default character set ("@%#*+=:~-.  ")
func DefaultOptions() Options {
	return Options{
		PixelRatio: DefaultPixelRatio(),
		Chars:      DefaultChars(),
	}
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
}
