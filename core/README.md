# ASCII Image Generator Core

Core package for converting images to ASCII art in Go. Provides flexible configuration options for pixel-to-character mapping and output customization.

## Features

- Configurable pixel-to-character ratio
- Customizable character sets
- Context-aware processing

## Usage

### Basic Example

```go
package main

import (
	"context"
	"image"
	"os"

	"github.com/fandasy/ASCIIimage/v2/core"
)

func main() {
	// Create default generation options
	opts := core.DefaultOptions()

	// Load your image (implement your own loader)
	img := loadImage("input.jpg")

	// Generate ASCII art image
	asciiImg, err := core.GenerateASCIIImage(context.Background(), img, opts)
	if err != nil {
		panic(err)
	}

	// Save the result (implement your own saver)
	saveImage("output.png", asciiImg)
}
```

### Advanced Configuration

```go
// Create custom generation options
opts := &core.Options{
    PixelRatio: core.PixelRatio{X: 2, Y: 3}, // 2x3 pixels → 1 ASCII char
    Chars: core.NewChars("01"),              // Custom character set
    Color: core.DefaultColor(),
}

// Applying a customization over a ready-made option
opts.WithFaceColor(color.RGBA{R: 122, G: 122, B: 122})

// Applying a setting to the default option
// opts := core.DefaultOptions().WithFaceColor(color.RGBA{R: 122, G: 122, B: 122})

asciiImg, err := core.GenerateASCIIImage(context.Background(), img, opts)
```

## API Reference

### Function

```go
// GenerateASCIIImage converts an image to ASCII art
func GenerateASCIIImage(ctx context.Context, img image.Image, opts_ptr *Options) (*image.RGBA, error)
```

### Options

```go
type Options struct {
    // PixelRatio defines how many original pixels map to one ASCII character
    PixelRatio PixelRatio // {X, Y}
    
    // Chars defines the character set to use (dark to light)
    Chars *Chars

    // Color specifies the foreground and background color scheme
    Color Color
}

// PixelRatio defines the pixel-to-character ratio
type PixelRatio struct {
    X, Y int
}

// Color represents a pair of foreground and background colors for ASCII art rendering.
// It ensures proper contrast between text and background.
type Color struct {
    // Face is the foreground/text color
    Face color.Color
    
    // Background is the canvas/background color
    Background color.Color
}
```

### Character Sets

```go
// Chars represents a character set mapping
type Chars [256]byte

// NewChars creates a new character set from a string (dark to light)
func NewChars(chars string) *Chars

// DefaultChars returns the default character set (@%#*+=:~-. )
func DefaultChars() *Chars
```
