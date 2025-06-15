package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/fandasy/ASCIIimage/v2/core"
	"github.com/fandasy/ASCIIimage/v2/validate"
	"golang.org/x/image/webp"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
)

// Client handles ASCII art generation with configurable options.
// It maintains an HTTP client for web requests and default options.
type Client struct {
	httpClient  *http.Client
	defaultOpts Options
}

// NewClient creates a new ASCII art client with custom HTTP client and options.
// If http.Client is nil, http.DefaultClient will be used.
// Options are validated before being set.
func NewClient(client *http.Client, opts_ptr *Options) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	opts := *opts_ptr

	opts.validate()

	return &Client{
		httpClient:  client,
		defaultOpts: opts,
	}
}

// NewDefaultClient creates a client with default HTTP client and options.
func NewDefaultClient() *Client {
	return &Client{
		httpClient:  http.DefaultClient,
		defaultOpts: *DefaultOptions(),
	}
}

var (
	// ErrFileNotFound indicates the requested file doesn't exist
	ErrFileNotFound = errors.New("file not found")

	// ErrPageNotFound indicates the requested URL returned 404
	ErrPageNotFound = errors.New("page not found")

	// ErrIncorrectFormat indicates unsupported image format
	ErrIncorrectFormat = errors.New("incorrect format")

	// ErrIncorrectUrl indicates malformed URL
	ErrIncorrectUrl = errors.New("incorrect url")
)

// GetFromFile reads an image from a file and converts it to ASCII art.
// Supported formats: PNG, JPEG, WebP.
//
// Parameters:
//   - ctx: Context for cancellation
//   - path: Path to image file
//   - opts: Optional conversion settings
//
// Returns:
//   - *image.RGBA: ASCII art image
//   - error: Possible errors:
//   - ErrFileNotFound
//   - ErrIncorrectFormat
//   - Other file operation or decoding errors
func (c *Client) GetFromFile(ctx context.Context, path string, opts ...Option) (*image.RGBA, error) {
	ext := filepath.Ext(path)
	if !validate.ContentType(ext, ".png", ".jpg", ".jpeg", ".webp") {
		return nil, fmt.Errorf("%w: %s", ErrIncorrectFormat, ext)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, ErrFileNotFound
	}

	defer file.Close()

	var img image.Image

	switch ext {
	case ".png":
		img, err = png.Decode(file)
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".webp":
		img, err = webp.Decode(file)
	}

	if err != nil {
		return nil, fmt.Errorf("img decoding failed: %w", err)
	}

	return c.GetFromImage(ctx, img, opts...)
}

// GetFromWebsite downloads an image from URL and converts it to ASCII art.
// Supported formats: PNG, JPEG, WebP.
//
// Parameters:
//   - ctx: Context for cancellation
//   - url: Image URL
//   - opts: Optional conversion settings
//
// Returns:
//   - *image.RGBA: ASCII art image
//   - error: Possible errors:
//   - ErrIncorrectUrl
//   - ErrPageNotFound
//   - ErrIncorrectFormat
//   - Other network or decoding errors
func (c *Client) GetFromWebsite(ctx context.Context, url string, opts ...Option) (*image.RGBA, error) {
	if !validate.URL(url) {
		return nil, ErrIncorrectUrl
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request failed: %w", err)
	}

	req.Close = true

	resp, err := c.httpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrPageNotFound
	}

	contentType := resp.Header.Get("Content-Type")

	if !validate.ContentType(contentType, "image/png", "image/jpeg", "image/webp") {
		return nil, fmt.Errorf("%w: %s", ErrIncorrectFormat, contentType)
	}

	var img image.Image

	switch contentType {
	case "image/png":
		img, err = png.Decode(resp.Body)
	case "image/jpeg":
		img, err = jpeg.Decode(resp.Body)
	case "image/webp":
		img, err = webp.Decode(resp.Body)
	}

	if err != nil {
		return nil, fmt.Errorf("img decoding failed: %w", err)
	}

	return c.GetFromImage(ctx, img, opts...)
}

// GetFromImage converts an existing image.Image to ASCII art.
//
// Parameters:
//   - ctx: Context for cancellation
//   - img: Source image
//   - opts: Optional conversion settings
//
// Returns:
//   - *image.RGBA: ASCII art image
//   - error: Context cancellation or processing errors
func (c *Client) GetFromImage(ctx context.Context, img image.Image, opts ...Option) (*image.RGBA, error) {

	ptrOpts := &c.defaultOpts

	if len(opts) != 0 {
		copyOpts := c.defaultOpts

		for _, opt := range opts {
			opt(&copyOpts)
		}

		ptrOpts = &copyOpts
	}

	ptrOpts.applyResizeOptions(img)

	return core.GenerateASCIIImage(ctx, img, &ptrOpts.Options)
}
