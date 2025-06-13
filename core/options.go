package core

const (
	defaultXPixelRatio = 1
	defaultYPixelRatio = 1
)

type PixelRatio struct {
	X, Y int
}

func DefaultPixelRatio() PixelRatio {
	return PixelRatio{defaultXPixelRatio, defaultYPixelRatio}
}

type Options struct {
	// x, y original pixels -> 1 ASCII char
	PixelRatio PixelRatio

	// chars for img generation
	Chars *Chars
}

func DefaultOptions() Options {
	return Options{
		PixelRatio: DefaultPixelRatio(),
		Chars:      DefaultChars(),
	}
}

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

type Option func(o *Options)

func WithPixelRatio(x, y int) Option {
	if x <= 0 {
		x = defaultXPixelRatio
	}
	if y <= 0 {
		y = defaultYPixelRatio
	}

	return func(o *Options) {
		o.PixelRatio = PixelRatio{x, y}
	}
}

func WithChars(c *Chars) Option {
	if c == nil {
		c = DefaultChars()
	}

	return func(o *Options) {
		o.Chars = c
	}
}
