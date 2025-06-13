package api

import (
	"github.com/fandasy/ASCIIimage/resize"
	"image"
)

// --------------------------------------------------
// There are default values for these parameters:
//
// Compress  = 0
// MaxWidth  = 10000 -> 100000px
// MaxHeight = 10000 -> 100000px
//
// --------------------------------------------------
// Default values can be activated by specifying:
//
// Compress  c < 0 || c > 99
// MaxWidth  <= 0
// MaxHeight <= 0
// --------------------------------------------------

const (
	// 1 	 = 10px
	// 10000 = 100000px

	defaultMaxWidth  uint = 10000
	defaultMaxHeight uint = 10000
)

// Options
//
// - Compress only in the range from 0 to 99
//
// - MaxWidth and MaxHeight are defined in the ratio 1 = 10px
type Options struct {
	Compress  uint8
	MaxWidth  uint
	MaxHeight uint
}

func DefaultOptions() Options {
	return Options{
		Compress:  0,
		MaxWidth:  defaultMaxWidth,
		MaxHeight: defaultMaxHeight,
	}
}

type Option func(*Options)

func WithCompress(compress uint8) Option {
	if compress < 0 || compress > 99 {
		compress = 0
	}

	return func(opts *Options) {
		opts.Compress = compress
	}
}

func WithMaxWidth(maxWidth uint) Option {
	if maxWidth <= 0 {
		maxWidth = defaultMaxWidth
	}

	return func(opts *Options) {
		opts.MaxWidth = maxWidth
	}
}

func WithMaxHeight(maxHeight uint) Option {
	if maxHeight <= 0 {
		maxHeight = defaultMaxHeight
	}

	return func(opts *Options) {
		opts.MaxHeight = maxHeight
	}
}

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
