package api

import (
	"context"
	"errors"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fandasy/ASCIIimage/v2/api"
)

func TestGetFromFile(t *testing.T) {
	const (
		validPath_1 = "../image/valid-img-1.jpg"
		validPath_2 = "../image/valid-img-2.jpg"
	)

	// Setup test files
	createTestImage := func(path string) {
		f, _ := os.Create(path)
		defer f.Close()
		img := image.NewRGBA(image.Rect(0, 0, 10, 10))
		png.Encode(f, img)
	}

	tests := []struct {
		name        string
		filePath    string
		createFile  bool
		expectError error
	}{
		{
			name:        "valid png file",
			filePath:    "test.png",
			createFile:  true,
			expectError: nil,
		},
		{
			name:        "first valid jpg file",
			filePath:    validPath_1,
			createFile:  false,
			expectError: nil,
		},
		{
			name:        "second valid jpg file",
			filePath:    validPath_2,
			createFile:  false,
			expectError: nil,
		},
		{
			name:        "non-existent file",
			filePath:    "nonexistent.png",
			createFile:  false,
			expectError: api.ErrFileNotFound,
		},
		{
			name:        "unsupported format",
			filePath:    "test.txt",
			createFile:  true,
			expectError: api.ErrIncorrectFormat,
		},
	}

	client := api.NewDefaultClient()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createFile {
				createTestImage(tt.filePath)
				defer os.Remove(tt.filePath)
			}

			_, err := client.GetFromFile(context.Background(), tt.filePath)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("GetFromFile() error = %v, want %v", err, tt.expectError)
			}
		})
	}
}

func TestGetFromWebsite(t *testing.T) {
	const (
		validUrl_1 = "https://www.gstatic.com/webp/gallery/4.jpg"
		validUrl_2 = "https://developers.google.com/static/search/docs/images/structured-data-in-image-results.png"
		validUrl_3 = "https://www.gstatic.com/webp/gallery/1.sm.webp"
	)

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/valid.png":
			w.Header().Set("Content-Type", "image/png")
			img := image.NewRGBA(image.Rect(0, 0, 10, 10))
			png.Encode(w, img)
		case "/404":
			w.WriteHeader(http.StatusNotFound)
		case "/wrong-content-type":
			w.Header().Set("Content-Type", "text/plain")
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	tests := []struct {
		name        string
		url         string
		expectError error
	}{
		{
			name:        "valid image url",
			url:         ts.URL + "/valid.png",
			expectError: nil,
		},
		{
			name:        "valid jpg url",
			url:         validUrl_1,
			expectError: nil,
		},
		{
			name:        "valid png url",
			url:         validUrl_2,
			expectError: nil,
		},
		{
			name:        "valid webp url",
			url:         validUrl_3,
			expectError: nil,
		},
		{
			name:        "non-existent page",
			url:         ts.URL + "/404",
			expectError: api.ErrPageNotFound,
		},
		{
			name:        "invalid url",
			url:         "invalid-url",
			expectError: api.ErrIncorrectUrl,
		},
		{
			name:        "wrong content type",
			url:         ts.URL + "/wrong-content-type",
			expectError: api.ErrIncorrectFormat,
		},
	}

	client := api.NewDefaultClient()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.GetFromWebsite(context.Background(), tt.url)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("GetFromWebsite() error = %v, want %v", err, tt.expectError)
			}
		})
	}
}

func TestGetFromImage(t *testing.T) {
	tests := []struct {
		name        string
		img         image.Image
		opts        []api.Option
		expectError bool
	}{
		{
			name:        "valid image with default options",
			img:         image.NewRGBA(image.Rect(0, 0, 10, 10)),
			opts:        nil,
			expectError: false,
		},
		{
			name:        "valid image with custom options",
			img:         image.NewRGBA(image.Rect(0, 0, 10, 10)),
			opts:        []api.Option{api.WithCompress(50)},
			expectError: false,
		},
	}

	client := api.NewDefaultClient()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.GetFromImage(context.Background(), tt.img, tt.opts...)
			if (err != nil) != tt.expectError {
				t.Errorf("GetFromImage() error = %v, want error %v", err, tt.expectError)
			}
		})
	}
}
