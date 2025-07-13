package core

import (
	"context"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"

	drawgray "github.com/fandasy/ASCIIimage/v2/pkg/draw-gray"
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

// GenerateASCIIImage converts an image to ASCII art.
// The conversion can be canceled using the provided context.
//
// Parameters:
//   - ctx: Context for cancellation
//   - img: Source image to convert
//   - opts: Conversion options (character set, pixel ratio, color)
//
// Returns:
//   - image.Image: Image containing the ASCII art
//   - error: Context cancellation error if operation was interrupted
func GenerateASCIIImage(ctx context.Context, img image.Image, opts_ptr *Options) (image.Image, error) {
	opts := *opts_ptr

	opts.validate()

	switch {
	case opts.Color.OriginalFace:
		// Drawing while preserving the original pixel color
		return generateASCIIImageWithOriginalColor(ctx, img, &opts)

	case opts.Color.TransparentBackground:
		// For a transparent background will need an alpha channel
		return generateASCIIImageToRGBA(ctx, img, &opts)

	case opts.Color.isGray():
		// image.Gray, image.Gray16
		return generateASCIIImageToGray(ctx, img, &opts)

	default:
		// image.RGBA, image.RGBA64, image.NRGBA, image.NRGBA64
		return generateASCIIImageToRGBA(ctx, img, &opts)
	}
}

// generateASCIIImageToRGBA returns image.RGBA, image.RGBA64, image.NRGBA, image.NRGBA64
func generateASCIIImageToRGBA(ctx context.Context, img image.Image, opts *Options) (image.Image, error) {
	bounds := img.Bounds()

	outputWidth := bounds.Max.X * (10 / opts.PixelRatio.X)
	outputHeight := bounds.Max.Y * (10 / opts.PixelRatio.Y)
	asciiImg := opts.Color.createDrawImage(outputWidth, outputHeight)

	lenAsciiLine := bounds.Max.X / opts.PixelRatio.X
	asciiLineBuf := make([]byte, 0, lenAsciiLine)

	if !opts.Color.TransparentBackground {
		draw.Draw(asciiImg, asciiImg.Bounds(), &image.Uniform{C: opts.Color.Background}, image.Point{}, draw.Src)
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y += opts.PixelRatio.Y {
		select {
		case <-ctx.Done():
			return asciiImg, ctx.Err()
		default:
		}

		asciiLine := asciiLineBuf[:0]

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

// generateASCIIImageToGray returns image.Gray, image.Gray16
func generateASCIIImageToGray(ctx context.Context, img image.Image, opts *Options) (image.Image, error) {
	bounds := img.Bounds()

	outputWidth := bounds.Max.X * (10 / opts.PixelRatio.X)
	outputHeight := bounds.Max.Y * (10 / opts.PixelRatio.Y)
	asciiImg := opts.Color.createDrawImage(outputWidth, outputHeight)

	lenAsciiLine := bounds.Max.X / opts.PixelRatio.X
	asciiLineBuf := make([]byte, 0, lenAsciiLine)

	// I don't check opts.Color.TransparentBackground because transparent background requires alpha channel

	drawgray.Draw(asciiImg, asciiImg.Bounds(), &image.Uniform{C: opts.Color.Background}, image.Point{})

	for y := bounds.Min.Y; y < bounds.Max.Y; y += opts.PixelRatio.Y {
		select {
		case <-ctx.Done():
			return asciiImg, ctx.Err()
		default:
		}

		asciiLine := asciiLineBuf[:0]

		for x := bounds.Min.X; x < bounds.Max.X; x += opts.PixelRatio.X {
			r, g, b, _ := img.At(x, y).RGBA()

			brightness := (r>>8 + g>>8 + b>>8) / 3

			asciiLine = append(asciiLine, opts.Chars[brightness])
		}

		scaledY := (y / opts.PixelRatio.Y) * 10

		point := fixed.Point26_6{X: fixed.I(0), Y: fixed.I(scaledY)}
		d := &drawgray.Drawer{
			Dst:  asciiImg,
			Src:  image.NewUniform(opts.Color.Face),
			Face: Face,
			Dot:  point,
		}
		d.DrawBytes(asciiLine)
	}

	return asciiImg, nil
}

// generateASCIIImageWithOriginalColor Drawing while preserving the original pixel color
func generateASCIIImageWithOriginalColor(ctx context.Context, img image.Image, opts *Options) (image.Image, error) {
	bounds := img.Bounds()

	outputWidth := bounds.Max.X * (10 / opts.PixelRatio.X)
	outputHeight := bounds.Max.Y * (10 / opts.PixelRatio.Y)
	asciiImg := opts.Color.createDrawImage(outputWidth, outputHeight)

	if !opts.Color.TransparentBackground {
		draw.Draw(asciiImg, asciiImg.Bounds(), &image.Uniform{C: opts.Color.Background}, image.Point{}, draw.Src)
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y += opts.PixelRatio.Y {
		select {
		case <-ctx.Done():
			return asciiImg, ctx.Err()
		default:
		}

		scaledY := (y / opts.PixelRatio.Y) * 10

		var (
			// RGBA
			xr, xg, xb, xa uint32

			scaledX int
			startX  int
			colorX  color.Color

			buf []byte
		)

		for x := bounds.Min.X; x < bounds.Max.X; x += opts.PixelRatio.X {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()

			brightness := (r>>8 + g>>8 + b>>8) / 3

			char := opts.Chars[brightness]

			// first initialization
			if colorX == nil {
				colorX = c
				buf = append(buf, char)
				startX = x
				continue
			}

			if r == xr && g == xg && b == xb && a == xa {
				buf = append(buf, char)
			} else {
				// draw previous segment
				scaledX = (startX / opts.PixelRatio.X) * 10
				d := &font.Drawer{
					Dst:  asciiImg,
					Src:  image.NewUniform(colorX),
					Face: Face,
					Dot:  fixed.Point26_6{X: fixed.I(scaledX), Y: fixed.I(scaledY)},
				}
				d.DrawBytes(buf)

				// reset
				colorX = c
				buf = []byte{char}
				startX = x
			}
		}

		// draw the remaining
		if len(buf) > 0 {
			scaledX = (startX / opts.PixelRatio.X) * 10
			d := &font.Drawer{
				Dst:  asciiImg,
				Src:  image.NewUniform(colorX),
				Face: Face,
				Dot:  fixed.Point26_6{X: fixed.I(scaledX), Y: fixed.I(scaledY)},
			}
			d.DrawBytes(buf)
		}
	}

	return asciiImg, nil
}
