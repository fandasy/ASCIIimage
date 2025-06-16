package core

import (
	"context"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
)

// Face provides a modified font for ASCII art rendering with:
//   - Character width: 10px
//   - Left padding: 2px
//   - Advance width: 10px
var Face = func() *basicfont.Face {
	copyFace := *basicfont.Face7x13
	face := &copyFace
	face.Width = 10
	face.Left = 2
	face.Advance = 10

	return face
}()

// GenerateASCIIImage converts an image to ASCII art rendered on an image.RGBA.
// The conversion can be canceled using the provided context.
//
// Parameters:
//   - ctx: Context for cancellation
//   - img: Source image to convert
//   - opts: Conversion options (character set, pixel ratio, color)
//
// Returns:
//   - *image.RGBA: Image containing the ASCII art
//   - error: Context cancellation error if operation was interrupted
func GenerateASCIIImage(ctx context.Context, img image.Image, opts_ptr *Options) (*image.RGBA, error) {
	opts := *opts_ptr

	opts.validate()

	bounds := img.Bounds()

	outputWidth := bounds.Max.X * (10 / opts.PixelRatio.X)
	outputHeight := bounds.Max.Y * (10 / opts.PixelRatio.X)

	asciiImg := image.NewRGBA(
		image.Rect(
			0,
			0,
			outputWidth,
			outputHeight,
		),
	)

	lenAsciiLine := bounds.Max.X / opts.PixelRatio.X

	draw.Draw(asciiImg, asciiImg.Bounds(), &image.Uniform{opts.Color.Background}, image.Point{}, draw.Src)

	for y := bounds.Min.Y; y < bounds.Max.Y; y += opts.PixelRatio.Y {
		select {
		case <-ctx.Done():
			return asciiImg, ctx.Err()
		default:
		}

		asciiLine := make([]byte, 0, lenAsciiLine)

		for x := bounds.Min.X; x < bounds.Max.X; x += opts.PixelRatio.X {
			r, g, b, _ := img.At(x, y).RGBA()

			brightness := (r>>8 + g>>8 + b>>8) / 3

			asciiLine = append(asciiLine, opts.Chars[brightness])
		}

		scaledY := (y / opts.PixelRatio.Y) * 10

		point := fixed.Point26_6{X: fixed.I(0), Y: fixed.I(scaledY)}
		d := &font.Drawer{
			Dst:  asciiImg,
			Src:  image.NewUniform(opts.Color.Face),
			Face: Face,
			Dot:  point,
		}
		d.DrawBytes(asciiLine)
	}

	return asciiImg, nil
}
