package api

import (
	"github.com/fandasy/ASCIIimage/v2/core"
	"github.com/fandasy/ASCIIimage/v2/pkg/resize"
	"image"
	"image/color"
)

const (
	// Default maximum dimensions (1 unit = 10px)
	defaultMaxWidth  uint = 10000 // 100,000px
	defaultMaxHeight uint = 10000 // 100,000px
)

// Options configure ASCII art generation behavior.
// Fields:
//   - Compress: Compression ratio (0-99)
//   - MaxWidth: Maximum width (1 = 10px)
//   - MaxHeight: Maximum height (1 = 10px)
//   - Core: core conversion options
type Options struct {
	Compress  uint8        // Image compression ratio (0-99)
	MaxWidth  uint         // Maximum width in 10px units
	MaxHeight uint         // Maximum height in 10px units
	Core      core.Options // Core ASCII generation options
}

// DefaultOptions returns the default configuration:
// - No compression
// - Maximum dimensions: 10000 (100,000px)
// - Default core options
func DefaultOptions() *Options {
	return &Options{
		Compress:  0,
		MaxWidth:  defaultMaxWidth,
		MaxHeight: defaultMaxHeight,
		Core:      *core.DefaultOptions(),
	}
}

func (o *Options) WithCompress(compress uint8) *Options {
	o.Compress = compress
	return o
}

func (o *Options) WithMaxWidth(maxWidth uint) *Options {
	o.MaxWidth = maxWidth
	return o
}

func (o *Options) WithMaxHeight(maxHeight uint) *Options {
	o.MaxHeight = maxHeight
	return o
}

func (o *Options) WithCoreOptions(options *core.Options) *Options {
	o.Core = *options
	return o
}

func (o *Options) WithPixelRatio(x, y int) *Options {
	o.Core.PixelRatio = core.PixelRatio{X: x, Y: y}
	return o
}

func (o *Options) WithChars(c *core.Chars) *Options {
	o.Core.Chars = c
	return o
}

func (o *Options) WithColor(color core.Color) *Options {
	o.Core.Color = color
	return o
}

func (o *Options) WithFaceColor(c color.Color) *Options {
	o.Core.Color.Face = c
	return o
}

func (o *Options) WithBackgroundColor(c color.Color) *Options {
	o.Core.Color.Background = c
	return o
}

// Option defines a function type for modifying Options
type Option func(*Options)

// WithCompress creates an Option to set compression ratio (0-99).
// Values outside 0-99 range will be clamped to 0.
func WithCompress(compress uint8) Option {
	if compress < 0 || compress > 99 {
		compress = 0
	}

	return func(opts *Options) {
		opts.Compress = compress
	}
}

// WithMaxWidth creates an Option to set maximum width.
// Values ≤ 0 will use defaultMaxWidth.
func WithMaxWidth(maxWidth uint) Option {
	if maxWidth <= 0 {
		maxWidth = defaultMaxWidth
	}

	return func(opts *Options) {
		opts.MaxWidth = maxWidth
	}
}

// WithMaxHeight creates an Option to set maximum height.
// Values ≤ 0 will use defaultMaxHeight.
func WithMaxHeight(maxHeight uint) Option {
	if maxHeight <= 0 {
		maxHeight = defaultMaxHeight
	}

	return func(opts *Options) {
		opts.MaxHeight = maxHeight
	}
}

// WithPixelRatio creates an Option to set pixel sampling ratio.
func WithPixelRatio(x, y int) Option {
	return func(opts *Options) {
		opts.Core.PixelRatio = core.PixelRatio{X: x, Y: y}
	}
}

// WithChars creates an Option to set custom character set.
func WithChars(c *core.Chars) Option {
	return func(opts *Options) {
		opts.Core.Chars = c
	}
}

// WithColor sets both foreground (face) and background colors for ASCII art generation.
//
// The color pair will be automatically validated to ensure proper contrast.
func WithColor(c core.Color) Option {
	return func(opts *Options) {
		opts.Core.Color = c
	}
}

// WithFaceColor sets only the foreground (text) color for ASCII art.
// The background color will remain unchanged unless explicitly set.
//
// The color will be automatically validated against the background:
//   - If nil, will use complementary color of background
//   - If same as background, will be adjusted for contrast
func WithFaceColor(c color.Color) Option {
	return func(opts *Options) {
		opts.Core.Color.Face = c
	}
}

// WithBackgroundColor sets only the background color for ASCII art.
// The foreground color will remain unchanged unless explicitly set.
//
// The color will be automatically validated against the foreground:
//   - If nil, will use complementary color of foreground
//   - If same as foreground, foreground will be adjusted for contrast
func WithBackgroundColor(c color.Color) Option {
	return func(opts *Options) {
		opts.Core.Color.Background = c
	}
}

// validate ensures option fields have valid values, setting defaults when needed.
func (o *Options) validate() {
	if o.Compress < 0 || o.Compress > 99 {
		o.Compress = 0
	}

	if o.MaxWidth <= 0 {
		o.MaxWidth = defaultMaxWidth
	}

	if o.MaxHeight <= 0 {
		o.MaxHeight = defaultMaxHeight
	}
}

// applyResizeOptions resizes the image according to options:
//   - Enforces MaxWidth / MaxHeight constraints
//   - Applies compression if specified.
//     Maintains aspect ratio during resizing
func (o *Options) applyResizeOptions(img image.Image) {
	bounds := img.Bounds()
	width := uint(bounds.Max.X)
	height := uint(bounds.Max.Y)

	resizeNeeded := width > o.MaxWidth || height > o.MaxHeight

	if o.Compress > 0 && o.Compress < 100 {
		resizeNeeded = true
	}

	if resizeNeeded {
		var newWidth, newHeight uint

		if width > o.MaxWidth || height > o.MaxHeight {
			// Maintain aspect ratio while clamping to max dimensions
			aspectRatio := float64(width) / float64(height)
			if width > o.MaxWidth {
				newWidth = o.MaxWidth
				newHeight = uint(float64(newWidth) / aspectRatio)
			}
			if newHeight > o.MaxHeight {
				newHeight = o.MaxHeight
				newWidth = uint(float64(newHeight) * aspectRatio)
			}
		}

		compressionFactor := uint(100 - o.Compress)
		newWidth = (width * (compressionFactor)) / 100
		newHeight = (height * (compressionFactor)) / 100

		img = resize.Resize(newWidth, newHeight, img)
	}
}
