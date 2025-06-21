# ASCII Image Generator

![Go Version](https://img.shields.io/badge/go-1.23+-blue.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Complete Go solution for converting images to ASCII art with flexible configuration options and multiple input sources.

## Features

- **Multiple input sources** (files, URLs, raw images)
- **Customizable output** (character sets, pixel ratios, colors)
- **Image processing** (resizing, compression)
- **Context-aware** operations

## Packages

| Package                  | Description |
|--------------------------|-------------|
| **[Core](core/)**        | Low-level ASCII generation logic |
| **[API](api/)**          | High-level client for common use cases |

## Installation

```bash
go get github.com/fandasy/ASCIIimage/v2
```

## Quick Start

### Basic Usage

```go
package main

import (
	"context"
	
	"github.com/fandasy/ASCIIimage/v2/api"
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
// Custom generation options
coreOpts := core.Options{
	PixelRatio: core.PixelRatio{X: 2, Y: 3},
	Chars:      core.NewChars("@%#*+=-:. "),
	Color:      core.DefaultColor(),
}

// Custom client
client := api.NewClient(
	&http.Client{Timeout: 30*time.Second},
	&api.Options{
		MaxWidth:  500, // 5000px
		Compress:  20,  // 20% compression
		Core:      coreOpts,
	},
)
```

## Examples

See example [file](example/file/main.go) and [url](example/url/main.go) directory of converted ascii images

Test [core](test/core/core_test.go), [api](test/api/api_test.go) - package of test keys

## Documentation

- [Core Package](core/README.md) - Low-level generation logic
- [API Package](api/README.md) - High-level client interface
- [Godoc Reference](https://pkg.go.dev/github.com/fandasy/ASCIIimage/v2) - Package: [Core](https://pkg.go.dev/github.com/fandasy/ASCIIimage/v2/core), [API](https://pkg.go.dev/github.com/fandasy/ASCIIimage/v2/api)

## License

MIT License. See [LICENSE](LICENSE) for details.
