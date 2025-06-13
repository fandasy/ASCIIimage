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

	"github.com/fandasy/ASCIIimage/core"
)

func main() {
	// Create default generator
	generator := core.DefaultGenerator()

	// Load your image (implement your own loader)
	img := loadImage("input.jpg")

	// Generate ASCII art image
	asciiImg, err := generator.GenerateASCIIImage(context.Background(), img)
	if err != nil {
		panic(err)
	}

	// Save the result (implement your own saver)
	saveImage("output.png", asciiImg)
}
```

### Advanced Configuration

```go
// Create custom generator with options
generator := core.NewGenerator(core.Options{
    PixelRatio: core.PixelRatio{X: 2, Y: 3}, // 2x3 pixels → 1 ASCII char
    Chars: core.NewChars("01"),              // Custom character set
})

// Generate with runtime options
asciiImg, err := generator.GenerateASCIIImage(
    ctx,
    img,
    core.WithPixelRatio(3, 3), // Override ratio for this generation
)
```

## API Reference

### Generator

```go
type Generator struct {
    // contains filtered or unexported fields
}

// NewGenerator creates a new generator with custom options
func NewGenerator(opts Options) *Generator

// DefaultGenerator creates a generator with default options
func DefaultGenerator() *Generator

// GenerateASCIIImage converts an image to ASCII art
func (g *Generator) GenerateASCIIImage(ctx context.Context, img image.Image, opts ...Option) (*image.RGBA, error)
```

### Options

```go
type Options struct {
    // PixelRatio defines how many original pixels map to one ASCII character
    PixelRatio PixelRatio // {X, Y}
    
    // Chars defines the character set to use (dark to light)
    Chars *Chars
}

// PixelRatio defines the pixel-to-character ratio
type PixelRatio struct {
    X, Y int
}

// Option allows runtime modification of options
type Option func(*Options)
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

## Configuration Options

### Pixel Ratio

Control how many original pixels map to one ASCII character:

```go
// 1x1 pixel → 1 character (highest detail)
WithPixelRatio(1, 1)

// 25x25 pixels → 1 character (medium detail)
WithPixelRatio(25, 25)

// 23/52 pixels → 1 character (custom)
WithPixelRatio(23, 52)
```

### Character Sets

Customize the ASCII characters used (from darkest to lightest):

```go
// Simple binary style
WithChars(NewChars("01"))

// Standard ASCII art
WithChars(NewChars("@%#*+=:~-.  "))
WithChars(DefaultChars())

// Custom gradient
WithChars(NewChars(" .:-=+*#%@"))
```
