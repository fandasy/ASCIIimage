package test

import (
	"context"
	"errors"
	"github.com/fandasy/ASCIIimage/api"
	"image/jpeg"
	"os"
	"testing"
)

func TestGetFromWebsite(t *testing.T) {
	const (
		validUrl_1 = "https://img1.akspic.ru/previews/1/6/0/6/7/176061/176061-yablochnyj_pejzazh-yabloko-illustracia-prirodnyj_landshaft-purpur-500x.jpg"
		validUrl_2 = "https://www.youloveit.ru/uploads/gallery/main/162/pikachu.png"
		validUrl_3 = "https://savvy.co.il/wp-content/themes/thesis/images/4.webp"

		exampleMaxWidth  = 3000 // 30000px
		exampleMaxHeight = 3000 // 30000px

		exampleChars = "0987654321"
	)

	type args struct {
		ctx  context.Context
		url  string
		opts []api.Option
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
				ctx: context.TODO(),
				url: validUrl_1,
				opts: []api.Option{
					api.WithCompress(50),
					api.WithMaxWidth(exampleMaxWidth),
					api.WithMaxHeight(exampleMaxHeight),
				},
			},
			wantErr: false,
		},
		{
			name: "Valid URL with 0% reduction",
			args: args{
				ctx: context.TODO(),
				url: validUrl_2,
				opts: []api.Option{
					api.WithMaxWidth(exampleMaxWidth),
					api.WithMaxHeight(exampleMaxHeight),
				},
			},
			wantErr: false,
		},
		{
			name: "Valid URL, webp format",
			args: args{
				ctx: context.TODO(),
				url: validUrl_3,
				opts: []api.Option{
					api.WithMaxWidth(exampleMaxWidth),
					api.WithMaxHeight(exampleMaxHeight),
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid URL",
			args: args{
				ctx: context.TODO(),
				url: "https://example.com/notValidURL",
				opts: []api.Option{
					api.WithMaxWidth(exampleMaxWidth),
					api.WithMaxHeight(exampleMaxHeight),
				},
			},
			wantErr: true,
			errStr:  api.ErrPageNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt

		client := api.NewDefaultClient()

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			img, err := client.GetFromWebsite(tt.args.ctx, tt.args.url, tt.args.opts...)

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
