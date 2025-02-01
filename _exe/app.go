package main

import (
	"context"
	"fmt"
	asciiimage "github.com/fandasy/ASCIIimage"
	"image"
	"image/jpeg"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	fmt.Println("Welcome to the ASCII Art Generator!")
	fmt.Println("Choose an option:")
	fmt.Println("1. Process an image from a file")
	fmt.Println("2. Process an image from a URL")
	fmt.Print("Enter your choice (1 or 2): ")

	var choice int
	_, err := fmt.Scan(&choice)
	if err != nil || (choice != 1 && choice != 2) {
		log.Fatal("Invalid choice. Please enter 1 or 2.")
	}

	switch choice {
	case 1:
		processFromFile()
	case 2:
		processFromWebsite()
	default:
		log.Fatal("Invalid choice. Exiting.")
	}
}

// processFromFile handles the logic for processing an image from a file.
func processFromFile() {
	fmt.Print("Enter the path to the image file: ")
	var path string
	_, err := fmt.Scan(&path)
	if err != nil {
		log.Fatal("Invalid input. Please provide a valid file path.")
	}

	fmt.Print("Enter the compression percentage (0-99): ")
	var compressionPercentage uint8
	_, err = fmt.Scan(&compressionPercentage)
	if err != nil || compressionPercentage > 100 {
		log.Fatal("Invalid compression percentage. Please enter a value between 0 and 100.")
	}

	fmt.Print("Enter the maximum width (1 = 10px): ")
	var maxWidth uint
	_, err = fmt.Scan(&maxWidth)
	if err != nil {
		log.Fatal("Invalid width. Please enter a positive integer.")
	}

	fmt.Print("Enter the maximum height (1 = 10px): ")
	var maxHeight uint
	_, err = fmt.Scan(&maxHeight)
	if err != nil {
		log.Fatal("Invalid height. Please enter a positive integer.")
	}

	opts := asciiimage.Options{
		Compress:  compressionPercentage,
		MaxWidth:  maxWidth,
		MaxHeight: maxHeight,
	}

	// Process the image
	ctx := context.Background()
	asciiImg, err := asciiimage.GetFromFile(ctx, path, opts)
	if err != nil {
		log.Fatalf("Error processing image: %v", err)
	}

	// Save the output
	outputPath := randomName() + ".png"
	saveImage(outputPath, asciiImg)
	fmt.Printf("ASCII art saved to %s\n", outputPath)
}

// processFromWebsite handles the logic for processing an image from a URL.
func processFromWebsite() {
	fmt.Print("Enter the URL of the image: ")
	var url string
	_, err := fmt.Scan(&url)
	if err != nil {
		log.Fatal("Invalid input. Please provide a valid URL.")
	}

	fmt.Print("Enter the compression percentage (0-99): ")
	var compressionPercentage uint8
	_, err = fmt.Scan(&compressionPercentage)
	if err != nil || compressionPercentage > 100 {
		log.Fatal("Invalid compression percentage. Please enter a value between 0 and 100.")
	}

	fmt.Print("Enter the maximum width (1 = 10px): ")
	var maxWidth uint
	_, err = fmt.Scan(&maxWidth)
	if err != nil {
		log.Fatal("Invalid width. Please enter a positive integer.")
	}

	fmt.Print("Enter the maximum height (1 = 10px): ")
	var maxHeight uint
	_, err = fmt.Scan(&maxHeight)
	if err != nil {
		log.Fatal("Invalid height. Please enter a positive integer.")
	}

	opts := asciiimage.Options{
		Compress:  compressionPercentage,
		MaxWidth:  maxWidth,
		MaxHeight: maxHeight,
	}

	// Process the image
	ctx := context.Background()
	asciiImg, err := asciiimage.GetFromWebsite(ctx, url, opts)
	if err != nil {
		log.Fatalf("Error processing image: %v", err)
	}

	// Save the output
	outputPath := randomName() + ".png"
	saveImage(outputPath, asciiImg)
	fmt.Printf("ASCII art saved to %s\n", outputPath)
}

// saveImage saves an image to a file.
func saveImage(path string, img *image.RGBA) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer file.Close()

	if err := jpeg.Encode(file, img, nil); err != nil {
		log.Fatalf("Error saving image: %v", err)
	}
}

func randomName() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	keyLength := 10

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[r.Intn(len(charset)-1)]
	}

	return string(shortKey)
}
