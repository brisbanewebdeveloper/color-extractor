package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"

	color_extractor "github.com/brisbanewebdeveloper/color-extractor"
)

func main() {
	file := "Fotolia_45549559_320_480.jpg"
	imageFile, _ := os.Open(file)
	defer imageFile.Close()

	image, _, _ := image.Decode(imageFile)
	colors := color_extractor.ExtractColors(image)

	fmt.Println(colors)
}
