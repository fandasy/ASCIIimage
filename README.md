# ASCII Image Generator

![Go Version](https://img.shields.io/badge/go-1.23+-blue.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Complete Go solution for converting images to ASCII art with flexible configuration options and multiple input sources.

## Features

- **Multiple input sources** (files, URLs, raw images)
- **Customizable output** (character sets, pixel ratios)
- **Image processing** (resizing, compression)
- **Context-aware** operations

## Packages

| Package | Description |
|---------|-------------|
| **[Core](core/README.md)** | Low-level ASCII generation logic |
| **[API](api/README.md)** | High-level client for common use cases |

## Installation

```bash
go get github.com/fandasy/ASCIIimage
```

## Quick Start

### Basic Usage

```go
package main

import (
	"context"
	"os"
	
	"github.com/fandasy/ASCIIimage/api"
)

func main() {
	client := api.NewDefaultClient()
	
	// From file
	asciiImg, _ := client.GetFromFile(context.Background(), "input.jpg")
	
	// From URL
	asciiImg, _ = client.GetFromWebsite(context.Background(), "https://example.com/image.png")
	
	// Save result
	saveImage("output.png", asciiImg)
}
```

### Advanced Configuration

```go
// Custom generator
generator := core.NewGenerator(core.Options{
	PixelRatio: core.PixelRatio{X: 2, Y: 3},
	Chars:      core.NewChars("@%#*+=-:. "),
})

// Custom client
client := api.NewClient(
	&http.Client{Timeout: 30*time.Second},
	generator,
	api.Options{
		MaxWidth:  500, // 5000px
		Compress:  20,  // 20% compression
	},
)
```

## Examples

See [example](example/) directory of converted ascii images

[TEST](test/) - package of usage keys

## Documentation

- [Core Package](core/README.md) - Low-level generation logic
- [API Package](api/README.md) - High-level client interface
- [Godoc Reference](https://pkg.go.dev/github.com/fandasy/ASCIIimage)

## License

MIT License. See [LICENSE](LICENSE) for details.