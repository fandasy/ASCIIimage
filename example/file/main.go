package main

import (
	"context"
	"fmt"
	"github.com/fandasy/ASCIIimage/v2/api"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	filepath_1 = "example/image/valid-img-1.jpg"
	filepath_2 = "example/image/valid-img-2.jpg"

	filename_1 = "save-1.png"
	filename_2 = "save-2.png"

	outputDir = "example/file"
)

func main() {
	httpClient := http.DefaultClient

	opts := api.DefaultOptions()

	client := api.NewClient(
		httpClient,
		opts,
	)

	img_1, err := client.GetFromFile(context.TODO(), filepath_1)
	if err != nil {
		log.Fatal(fmt.Errorf("GetFromFile: %v", err))
	}

	path_1 := filepath.Join(outputDir, filename_1)

	if err := saveImage(path_1, img_1); err != nil {
		log.Fatal(fmt.Errorf("saveImage: %v", err))
	}

	img_2, err := client.GetFromFile(context.TODO(), filepath_2)
	if err != nil {
		log.Fatal(fmt.Errorf("GetFromFile: %v", err))
	}

	path_2 := filepath.Join(outputDir, filename_2)

	if err := saveImage(path_2, img_2); err != nil {
		log.Fatal(fmt.Errorf("saveImage: %v", err))
	}

	log.Println("example/file case successfully complied")
}

func saveImage(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", path, err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode image to JPEG: %v", err)
	}

	log.Printf("saved image to file %s", path)

	return nil
}
