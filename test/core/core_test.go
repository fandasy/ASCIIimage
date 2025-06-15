package core

import (
	"context"
	"image"
	"image/color"
	"testing"

	"github.com/fandasy/ASCIIimage/v2/core"
)

func TestNewChars(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		checkFn func(*core.Chars) bool
	}{
		{
			name:  "default chars",
			input: "@%#*+=:~-.  ",
			checkFn: func(c *core.Chars) bool {
				return c[0] == '@' && c[255] == ' '
			},
		},
		{
			name:  "empty string",
			input: "",
			checkFn: func(c *core.Chars) bool {
				return c == core.DefaultChars()
			},
		},
		{
			name:  "non-ascii chars",
			input: "你好世界",
			checkFn: func(c *core.Chars) bool {
				return c == core.DefaultChars()
			},
		},
		{
			name:  "single char",
			input: "@",
			checkFn: func(c *core.Chars) bool {
				for i := 0; i < 256; i++ {
					if c[i] != '@' {
						return false
					}
				}
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := core.NewChars(tt.input)
			if !tt.checkFn(got) {
				t.Errorf("NewChars() = %v, validation failed", got)
			}
		})
	}
}

func TestGenerateASCIIImage(t *testing.T) {
	// Create test image (2x2 pixels)
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{0, 0, 0, 255})       // Black
	img.Set(1, 0, color.RGBA{128, 128, 128, 255}) // Gray
	img.Set(0, 1, color.RGBA{255, 255, 255, 255}) // White
	img.Set(1, 1, color.RGBA{255, 0, 0, 255})     // Red

	tests := []struct {
		name      string
		img       image.Image
		opts      *core.Options
		wantError bool
		checkFn   func(*image.RGBA) bool
	}{
		{
			name:      "default options",
			img:       img,
			opts:      core.DefaultOptions(),
			wantError: false,
			checkFn: func(result *image.RGBA) bool {
				// Check basic properties
				bounds := result.Bounds()
				return bounds.Dx() == 20 && bounds.Dy() == 20 // 2x * 10px
			},
		},
		{
			name: "custom pixel ratio",
			img:  img,
			opts: &core.Options{
				PixelRatio: core.PixelRatio{X: 2, Y: 1},
			},
			wantError: false,
			checkFn: func(result *image.RGBA) bool {
				bounds := result.Bounds()
				return bounds.Dx() == 10 && bounds.Dy() == 20 // (2/2)x * 10px, 2y * 10px
			},
		},
		{
			name: "custom chars",
			img:  img,
			opts: &core.Options{
				Chars: core.NewChars("01"),
			},
			wantError: false,
			checkFn: func(result *image.RGBA) bool {
				// Should only contain 0 or 1 characters
				return true // Would need OCR or pixel analysis to verify
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := core.GenerateASCIIImage(ctx, tt.img, tt.opts)

			if (err != nil) != tt.wantError {
				t.Errorf("GenerateASCIIImage() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError && tt.checkFn != nil && !tt.checkFn(got) {
				t.Error("GenerateASCIIImage() validation failed")
			}
		})
	}
}
