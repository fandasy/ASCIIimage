package asciiimage

import (
	"context"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
)

type Generator struct {
	defaultOpts Options
}

func NewGenerator(opts Options) *Generator {
	opts.validate()

	return &Generator{opts}
}

func DefaultGenerator() *Generator {
	return &Generator{
		defaultOpts: DefaultOptions(),
	}
}

func init() {
	copyFace := *basicfont.Face7x13
	Face = &copyFace
	Face.Width = 10
	Face.Left = 2
	Face.Advance = 10
}

var Face *basicfont.Face

// GenerateASCIIImage converts the image to ASCII art.
func (g *Generator) GenerateASCIIImage(ctx context.Context, img image.Image, opts ...Option) (*image.RGBA, error) {
	ptrOpts := &g.defaultOpts

	if len(opts) != 0 {
		copyOpts := g.defaultOpts

		for _, opt := range opts {
			opt(&copyOpts)
		}

		ptrOpts = &copyOpts
	}

	bounds := img.Bounds()
	asciiWidth := bounds.Max.X
	asciiHeight := bounds.Max.Y

	asciiImg := image.NewRGBA(image.Rect(0, 0, asciiWidth*10, asciiHeight*10))

	draw.Draw(asciiImg, asciiImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	for y := bounds.Min.Y; y < bounds.Max.Y; y += ptrOpts.PixelRatio.Y {
		select {
		case <-ctx.Done():
			return asciiImg, ctx.Err()
		default:
		}

		asciiLine := make([]byte, 0, bounds.Max.X)
		for x := bounds.Min.X; x < bounds.Max.X; x += ptrOpts.PixelRatio.X {
			r, g, b, _ := img.At(x, y).RGBA()

			brightness := (r>>8 + g>>8 + b>>8) / 3

			asciiLine = append(asciiLine, ptrOpts.Chars[brightness])
		}

		point := fixed.Point26_6{X: fixed.I(0), Y: fixed.I(y * 10)}
		d := &font.Drawer{
			Dst:  asciiImg,
			Src:  image.NewUniform(color.Black),
			Face: Face,
			Dot:  point,
		}
		d.DrawBytes(asciiLine)
	}

	return asciiImg, nil
}
