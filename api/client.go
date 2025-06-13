package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/fandasy/ASCIIimage/core"
	"github.com/fandasy/ASCIIimage/validate"
	"golang.org/x/image/webp"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
)

type Client struct {
	httpClient  *http.Client
	generator   *core.Generator
	defaultOpts Options
}

func NewClient(client *http.Client, generator *core.Generator, opts Options) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	if generator == nil {
		generator = core.DefaultGenerator()
	}

	opts.validate()

	return &Client{
		httpClient:  client,
		generator:   generator,
		defaultOpts: opts,
	}
}

func NewDefaultClient() *Client {
	return &Client{
		httpClient:  http.DefaultClient,
		generator:   core.DefaultGenerator(),
		defaultOpts: DefaultOptions(),
	}
}

var (
	ErrFileNotFound    = errors.New("file not found")
	ErrPageNotFound    = errors.New("page not found")
	ErrIncorrectFormat = errors.New("incorrect format")
	ErrIncorrectUrl    = errors.New("incorrect url")
)

// GetFromFile reads an image from a file and converts it to an ASCII art image.
//
// Possible output errors:
// ErrFileNotFound,
// ErrIncorrectFormat
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

// GetFromWebsite downloads an image from a URL and converts it to an ASCII art image.
//
// Possible output errors:
// ErrIncorrectUrl,
// ErrPageNotFound,
// ErrIncorrectFormat
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

// GetFromImage processes the image and generates ASCII art.
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

	return c.generator.GenerateASCIIImage(ctx, img)
}
