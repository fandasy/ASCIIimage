package test

import (
	"context"
	"errors"
	"github.com/fandasy/ASCIIimage/v2/api"
	"image/jpeg"
	"os"
	"testing"
)

func TestGetFromFile(t *testing.T) {
	const (
		validPath_1 = "../example/valid-img-1.jpg"
		validPath_2 = "../example/valid-img-2.jpg"

		exampleMaxWidth  = 3000 // 30000px
		exampleMaxHeight = 3000 // 30000px

		exampleChars = "0987654321"
	)

	type args struct {
		ctx  context.Context
		path string
		opts []api.Option
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
				ctx:  context.TODO(),
				path: validPath_1,
				opts: []api.Option{
					api.WithMaxWidth(exampleMaxWidth),
					api.WithMaxHeight(exampleMaxHeight),
				},
			},
			wantErr: false,
		},
		{
			name: "Valid path 2 with 0% reduction",
			args: args{
				ctx:  context.TODO(),
				path: validPath_2,
				opts: nil,
			},
			wantErr: false,
		},
		{
			name: "File not found",
			args: args{
				path: "test-image/example.png",
			},
			wantErr: true,
			errStr:  api.ErrFileNotFound,
		},
		{
			name: "Incorrect Format",
			args: args{
				path: "test-image/example.txt",
			},
			wantErr: true,
			errStr:  api.ErrIncorrectFormat,
		},
	}
	for _, tt := range tests {
		tt := tt

		client := api.NewDefaultClient()

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			img, err := client.GetFromFile(tt.args.ctx, tt.args.path, tt.args.opts...)

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
