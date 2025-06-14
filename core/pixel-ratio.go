package core

const (
	defaultXPixelRatio = 1
	defaultYPixelRatio = 1
)

// PixelRatio defines the sampling ratio between source image and ASCII output
type PixelRatio struct {
	X, Y int // Horizontal and vertical sampling ratios
}

// DefaultPixelRatio returns the default 1:1 pixel ratio
func DefaultPixelRatio() PixelRatio {
	return PixelRatio{defaultXPixelRatio, defaultYPixelRatio}
}
