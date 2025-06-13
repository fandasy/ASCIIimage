# ASCII Image API Client

High-level API client for converting images to ASCII art from various sources. Built on top of the `core` package with additional convenience methods for file and web image handling.

## Features

- Load images from local files
- Download images from URLs
- Automatic image format detection (PNG, JPEG, WebP)
- Configurable compression and resizing
- Context-aware operations
- Custom HTTP client support

## Usage

### Basic Example

```go
package main

import (
	"context"
	"os"
	
	"github.com/fandasy/ASCIIimage/api/v2"
)

func main() {
	// Create default client
	client := api.NewDefaultClient()

	// From local file
	asciiImgFromFile, err := client.GetFromFile(context.Background(), "input.jpg")
	if err != nil {
		panic(err)
	}

	// From URL
	asciiImgFromWebsite, err := client.GetFromWebsite(context.Background(), "https://example.com/image.png")
	if err != nil {
		panic(err)
	}

	// Save the result (implement your own saver)
	
	saveImage("ascii_1.png", asciiImgFromFile)
	
	saveImage("ascii_2.png", asciiImgFromWebsite)
}
```

### Advanced Configuration

```go
// Create custom HTTP client with timeout
httpClient := &http.Client{
	Timeout: 30 * time.Second,
}

// Create custom generator
generator := core.NewGenerator(core.Options{
	PixelRatio: core.PixelRatio{X: 2, Y: 3},
})

// Create API client with custom options
client := api.NewClient(
	httpClient,
	generator,
	api.Options{
		MaxWidth:  500,  // 5000px (1 unit = 10px)
		MaxHeight: 300,  // 3000px
		Compress:  30,   // 30% compression
	},
)

// Generate with additional runtime options
asciiImg, err := client.GetFromFile(
	ctx,
	"large-image.jpg",
	api.WithMaxWidth(200),  // Override max width for this request
)
```

## API Reference

### Client

```go
type Client struct {
	// contains filtered or unexported fields
}

// NewClient creates a new client with custom configuration
func NewClient(client *http.Client, generator *core.Generator, opts Options) *Client

// NewDefaultClient creates a client with default configuration
func NewDefaultClient() *Client

// GetFromFile reads an image from file and converts to ASCII art
func (c *Client) GetFromFile(ctx context.Context, path string, opts ...Option) (*image.RGBA, error)

// GetFromWebsite downloads an image from URL and converts to ASCII art
func (c *Client) GetFromWebsite(ctx context.Context, url string, opts ...Option) (*image.RGBA, error)

// GetFromImage converts an existing image to ASCII art
func (c *Client) GetFromImage(ctx context.Context, img image.Image, opts ...Option) (*image.RGBA, error)
```

### Options

```go
type Options struct {
	Compress  uint8 // Compression percentage (0-99)
	MaxWidth  uint  // Maximum width (1 unit = 10px)
	MaxHeight uint  // Maximum height (1 unit = 10px)
}

// Option allows runtime modification of options
type Option func(*Options)

// WithCompress sets compression level (0-99)
func WithCompress(compress uint8) Option

// WithMaxWidth sets maximum width (1 unit = 10px)
func WithMaxWidth(maxWidth uint) Option

// WithMaxHeight sets maximum height (1 unit = 10px)
func WithMaxHeight(maxHeight uint) Option
```

### Error Handling

The package defines several common errors:

```go
var (
	ErrFileNotFound    = errors.New("file not found")
	ErrPageNotFound    = errors.New("page not found")
	ErrIncorrectFormat = errors.New("incorrect format")
	ErrIncorrectUrl    = errors.New("incorrect url")
)
```

## Configuration Options

### Image Resizing

Control maximum dimensions and compression:

```go
// Set maximum dimensions (5000x3000 pixels)
WithMaxWidth(500)  // 500 * 10px = 5000px
WithMaxHeight(300) // 300 * 10px = 3000px

// Apply 25% compression
WithCompress(25)
```

### Custom HTTP Client

```go
// Create client with custom timeout and transport
httpClient := &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns: 10,
	},
}

client := api.NewClient(httpClient, generator, opts)
```

## Supported Image Formats

- PNG
- JPEG/JPG
- WebP

## Error Handling

All methods return errors that can be checked against the defined error variables:

```go
img, err := client.GetFromWebsite(ctx, url)
if errors.Is(err, api.ErrPageNotFound) {
	// Handle 404
} else if errors.Is(err, api.ErrIncorrectFormat) {
	// Handle unsupported format
}
```
