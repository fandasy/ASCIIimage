package asciiimage

import (
	"context"
	"errors"
	"image/jpeg"
	"os"
	"testing"
)

const (
	validPath_1 = "example/valid-img-1.jpg"
	validPath_2 = "example/valid-img-2.jpg"

	validUrl_1 = "https://img1.akspic.ru/previews/1/6/0/6/7/176061/176061-yablochnyj_pejzazh-yabloko-illustracia-prirodnyj_landshaft-purpur-500x.jpg"
	validUrl_2 = "https://www.youloveit.ru/uploads/gallery/main/162/pikachu.png"
	validUrl_3 = "https://savvy.co.il/wp-content/themes/thesis/images/4.webp"

	exampleMaxWidth  = 3000 // 30000px
	exampleMaxHeight = 3000 // 30000px

	chars = "@%#*+=:~-. "
)

// --------------------------------------------------
// There are default values for these parameters:
//
// Compress  = 0
// MaxWidth  = 10000 -> 100000px
// MaxHeight = 10000 -> 100000px
// Chars     = "@%#*+=:~-. "
//
// --------------------------------------------------
// Default values can be activated by specifying:
//
// Compress  >= 100
// MaxWidth  == 0
// MaxHeight == 0
// Chars     == ""
// --------------------------------------------------

func TestGetFromFile(t *testing.T) {
	type args struct {
		ctx       context.Context
		path      string
		compress  uint8
		maxWidth  uint
		maxHeight uint
		chars     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		errStr  error
	}{
		{
			name: "Valid path 1 with 0% reduction",
			args: args{
				ctx:       context.TODO(),
				path:      validPath_1,
				compress:  0,
				maxWidth:  exampleMaxWidth,
				maxHeight: exampleMaxHeight,
				chars:     chars,
			},
			wantErr: false,
		},
		{
			name: "Valid path 2 with 0% reduction",
			args: args{
				ctx:       context.TODO(),
				path:      validPath_2,
				compress:  0,
				maxWidth:  exampleMaxWidth,
				maxHeight: exampleMaxHeight,
				chars:     chars,
			},
			wantErr: false,
		},
		{
			name: "File not found",
			args: args{
				path: "test-image/example.png",
			},
			wantErr: true,
			errStr:  ErrFileNotFound,
		},
		{
			name: "Incorrect Format",
			args: args{
				path: "test-image/example.txt",
			},
			wantErr: true,
			errStr:  ErrIncorrectFormat,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			opts := Options{
				Compress:  tt.args.compress,
				MaxWidth:  tt.args.maxWidth,
				MaxHeight: tt.args.maxHeight,
				Chars:     tt.args.chars,
			}

			img, err := GetFromFile(tt.args.ctx, tt.args.path, opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {

				fileName := tt.name + ".jpg"
				file, err := os.Create(fileName)
				if err != nil {
					t.Errorf("failed to create file %s: %v", fileName, err)
					return
				}
				defer file.Close()

				if err := jpeg.Encode(file, img, nil); err != nil {
					t.Errorf("failed to encode image to JPEG: %v", err)
				}

			} else {

				if !errors.Is(err, tt.errStr) {
					t.Errorf("Get() error = %v, errStr %v", err, tt.errStr)
				}

			}
		})
	}
}

func TestGetFromWebsite(t *testing.T) {
	type args struct {
		ctx       context.Context
		url       string
		compress  uint8
		maxWidth  uint
		maxHeight uint
		chars     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		errStr  error
	}{
		{
			name: "Valid URL with 50% reduction",
			args: args{
				ctx:       context.TODO(),
				url:       validUrl_1,
				compress:  50,
				maxWidth:  exampleMaxWidth,
				maxHeight: exampleMaxHeight,
				chars:     chars,
			},
			wantErr: false,
		},
		{
			name: "Valid URL with 0% reduction",
			args: args{
				ctx:       context.TODO(),
				url:       validUrl_2,
				compress:  0,
				maxWidth:  exampleMaxWidth,
				maxHeight: exampleMaxHeight,
				chars:     chars,
			},
			wantErr: false,
		},
		{
			name: "Valid URL, webp format",
			args: args{
				ctx:       context.TODO(),
				url:       validUrl_3,
				compress:  0,
				maxWidth:  exampleMaxWidth,
				maxHeight: exampleMaxHeight,
				chars:     chars,
			},
			wantErr: false,
		},
		{
			name: "Invalid URL",
			args: args{
				ctx:       context.TODO(),
				url:       "https://example.com/notValidURL",
				compress:  0,
				maxWidth:  exampleMaxWidth,
				maxHeight: exampleMaxHeight,
				chars:     chars,
			},
			wantErr: true,
			errStr:  ErrPageNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			opts := Options{
				Compress:  tt.args.compress,
				MaxWidth:  tt.args.maxWidth,
				MaxHeight: tt.args.maxHeight,
				Chars:     tt.args.chars,
			}

			img, err := GetFromWebsite(tt.args.ctx, tt.args.url, opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {

				fileName := tt.name + ".jpg"
				file, err := os.Create(fileName)
				if err != nil {
					t.Errorf("failed to create file %s: %v", fileName, err)
					return
				}
				defer file.Close()

				if err := jpeg.Encode(file, img, nil); err != nil {
					t.Errorf("failed to encode image to JPEG: %v", err)
				}

			} else {

				if !errors.Is(err, tt.errStr) {
					t.Errorf("Get() error = %v, errStr %v", err, tt.errStr)
				}

			}
		})
	}
}
