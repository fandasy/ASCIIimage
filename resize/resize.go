package resize

import (
	"image"
	"image/color"
	"math"
)

// Resize resizes the image to the specified width and height using the Lanczos2 algorithm.
func Resize(newWidth, newHeight uint, img image.Image) image.Image {
	bounds := img.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	// Create a new image with the desired dimensions
	dst := image.NewRGBA(image.Rect(0, 0, int(newWidth), int(newHeight)))

	// Calculate scaling factors
	xScale := float64(srcWidth) / float64(newWidth)
	yScale := float64(srcHeight) / float64(newHeight)

	// Apply Lanczos2 resampling
	for y := 0; y < int(newHeight); y++ {
		for x := 0; x < int(newWidth); x++ {
			// Map destination coordinates to source coordinates
			srcX := (float64(x)+0.5)*xScale - 0.5
			srcY := (float64(y)+0.5)*yScale - 0.5

			// Get the interpolated color at the source coordinates
			c := lanczos2Interpolate(img, srcX, srcY)

			// Set the color in the destination image
			dst.Set(x, y, c)
		}
	}

	return dst
}

// lanczos2Interpolate performs Lanczos2 interpolation at the given coordinates.
func lanczos2Interpolate(img image.Image, x, y float64) color.Color {
	// Get the integer coordinates of the 4x4 grid around (x, y)
	x0 := int(math.Floor(x)) - 1
	y0 := int(math.Floor(y)) - 1

	var r, g, b, a float64

	// Iterate over the 4x4 grid
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			// Calculate the weight using the Lanczos2 kernel
			dx := x - float64(x0+i)
			dy := y - float64(y0+j)
			wx := lanczos2Kernel(dx)
			wy := lanczos2Kernel(dy)
			weight := wx * wy

			// Get the color at the current grid point
			c := img.At(x0+i, y0+j)
			cr, cg, cb, ca := c.RGBA()

			// Accumulate the weighted color components
			r += float64(cr) * weight
			g += float64(cg) * weight
			b += float64(cb) * weight
			a += float64(ca) * weight
		}
	}

	// Normalize the color components
	r = clampFloat(r, 0, 0xffff)
	g = clampFloat(g, 0, 0xffff)
	b = clampFloat(b, 0, 0xffff)
	a = clampFloat(a, 0, 0xffff)

	return color.RGBA64{
		R: uint16(r),
		G: uint16(g),
		B: uint16(b),
		A: uint16(a),
	}
}

// lanczos2Kernel computes the Lanczos2 kernel value.
func lanczos2Kernel(x float64) float64 {
	if x == 0 {
		return 1
	}
	if math.Abs(x) >= 2 {
		return 0
	}
	return (math.Sin(math.Pi*x) * math.Sin(math.Pi*x/2)) / (math.Pi * math.Pi * x * x / 2)
}

// clampFloat ensures a float value is within a specified range.
func clampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
