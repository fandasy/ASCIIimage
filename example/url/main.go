package main

import (
	"context"
	"fmt"
	"github.com/fandasy/ASCIIimage/v2/api"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
)

const (
	url_1 = "https://www.gstatic.com/webp/gallery/4.jpg"
	url_2 = "https://developers.google.com/static/search/docs/images/structured-data-in-image-results.png"
	url_3 = "https://www.gstatic.com/webp/gallery3/1_webp_ll.sm.png"

	filename_1 = "save-1.jpg"
	filename_2 = "save-2.jpg"
	filename_3 = "save-3.jpg"

	outputDir = "example/url"
)

func main() {
	client := api.NewDefaultClient()

	// start 1
	img_1, err := client.GetFromWebsite(context.TODO(), url_1)
	if err != nil {
		log.Fatal(fmt.Errorf("GetFromWebsite: %v", err))
	}

	path_1 := filepath.Join(outputDir, filename_1)

	if err := saveImage(path_1, img_1); err != nil {
		log.Fatal(fmt.Errorf("saveImage: %v", err))
	}
	// end

	// start 2
	img_2, err := client.GetFromWebsite(context.TODO(), url_2)
	if err != nil {
		log.Fatal(fmt.Errorf("GetFromWebsite: %v", err))
	}

	path_2 := filepath.Join(outputDir, filename_2)

	if err := saveImage(path_2, img_2); err != nil {
		log.Fatal(fmt.Errorf("saveImage: %v", err))
	}
	// end

	// start 3
	img_3, err := client.GetFromWebsite(context.TODO(), url_3)
	if err != nil {
		log.Fatal(fmt.Errorf("GetFromWebsite: %v", err))
	}

	path_3 := filepath.Join(outputDir, filename_3)

	if err := saveImage(path_3, img_3); err != nil {
		log.Fatal(fmt.Errorf("saveImage: %v", err))
	}
	// end

	log.Println("example/url case successfully complied")
}

func saveImage(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", path, err)
	}
	defer file.Close()

	if err := jpeg.Encode(file, img, nil); err != nil {
		return fmt.Errorf("failed to encode image to JPEG: %v", err)
	}

	log.Printf("saved image to file %s", path)

	return nil
}
