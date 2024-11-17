package asciiimage

import (
	"context"
	"errors"
	"image/jpeg"
	"os"
	"testing"
)

const (
	validPath_1 = "example/8a525c04a725d2d65bc5de411715f428.jpg"
	validPath_2 = "example/c80ece0becfb487ad70b0f379d4bfb78.jpg"

	validUrl_1 = "https://ir.ozone.ru/s3/multimedia-7/c1000/6755179327.jpg"
	validUrl_2 = "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRVNyUgiDjGP2BXtaLCP48USXG0l9sttdYNgw&s"
	validUrl_3 = "https://a.storyblok.com/f/178900/1000x563/8c02cd867d/2bcc3020d88ac2cfc03709aca149d5a21681300939_main.jpg/m/filters:quality(95)format(webp)"

	maxWidth  = 3000 // 30000px
	maxHeight = 1500 // 15000px

	chars = "@%#*+=:~-. "
)

// --------------------------------------------------
// There are default values for these parameters:
//
// reductionPercentage = 0.0
// maxWidth  = 5000 -> 50000px
// maxHeight = 5000 -> 50000px
// chars     = "@%#*+=:~-. "
//
// --------------------------------------------------
// Default values can be activated by specifying:
//
// reductionPercentage < 0 || reductionPercentage > 1
// maxWidth  <= 0
// maxHeight <= 0
// chars     == ""
// --------------------------------------------------

func TestGetFromFile(t *testing.T) {
	type args struct {
		path                string
		reductionPercentage float64
		maxWidth            int
		maxHeight           int
		chars               string
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
				path:                validPath_1,
				reductionPercentage: 0.0,
				maxWidth:            maxWidth,
				maxHeight:           maxHeight,
				chars:               chars,
			},
			wantErr: false,
		},
		{
			name: "Valid path 2 with 0% reduction",
			args: args{
				path:                validPath_2,
				reductionPercentage: 0.0,
				maxWidth:            maxWidth,
				maxHeight:           maxHeight,
				chars:               chars,
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

			img, err := GetFromFile(tt.args.path, tt.args.reductionPercentage, tt.args.maxWidth, tt.args.maxHeight, tt.args.chars)

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
		ctx                 context.Context
		url                 string
		reductionPercentage float64
		maxWidth            int
		maxHeight           int
		chars               string
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
				ctx:                 context.Background(),
				url:                 validUrl_1,
				reductionPercentage: 0.5,
				maxWidth:            maxWidth,
				maxHeight:           maxHeight,
				chars:               chars,
			},
			wantErr: false,
		},
		{
			name: "Valid URL with 0% reduction",
			args: args{
				ctx:                 context.Background(),
				url:                 validUrl_2,
				reductionPercentage: 0.0,
				maxWidth:            maxWidth,
				maxHeight:           maxHeight,
				chars:               chars,
			},
			wantErr: false,
		},
		{
			name: "Valid URL, webp format",
			args: args{
				ctx:                 context.Background(),
				url:                 validUrl_3,
				reductionPercentage: 0.0,
				maxWidth:            maxWidth,
				maxHeight:           maxHeight,
				chars:               chars,
			},
			wantErr: false,
		},
		{
			name: "Invalid URL",
			args: args{
				ctx:                 context.Background(),
				url:                 "https://example.com/notValidURL",
				reductionPercentage: 0.0,
				maxWidth:            maxWidth,
				maxHeight:           maxHeight,
				chars:               chars,
			},
			wantErr: true,
			errStr:  ErrPageNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			img, err := GetFromWebsite(tt.args.ctx, tt.args.url, tt.args.reductionPercentage, tt.args.maxWidth, tt.args.maxHeight, tt.args.chars)

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
