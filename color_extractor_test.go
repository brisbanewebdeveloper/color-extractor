package color_extractor

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"os"
	"testing"
)

func TestExtractColors(t *testing.T) {
	white := color.RGBA{R: 225, G: 255, B: 255, A: 255}
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	green := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	transparent := color.RGBA{R: 0, G: 0, B: 0, A: 0}
	semiTransparentRed := color.RGBA{R: 255, G: 0, B: 0, A: 127}

	testCases := map[string]struct {
		Image           image.Image
		ExtractedColors []Color
	}{
		"Empty file": {
			Image:           imageFromColors([]Color{}),
			ExtractedColors: []Color{},
		},
		"Single pixel": {
			Image: imageFromColors([]Color{
				{Color: red, Count: 0.4980392156862745},
			}),
			ExtractedColors: []Color{
				{Color: red, Count: 1},
			},
		},
		"One color": {
			Image: imageFromColors([]Color{
				{Color: white},
				{Color: white},
				{Color: white},
				{Color: white},
				{Color: white},
			}),
			ExtractedColors: []Color{
				{Color: white, Count: 5},
			},
		},
		"Transparent image": {
			Image: imageFromColors([]Color{
				{Color: white},
				{Color: white},
				{Color: white},
				{Color: transparent},
			}),
			ExtractedColors: []Color{
				{Color: white, Count: 3},
			},
		},
		"Semitransparent single pixel": {
			Image: imageFromColors([]Color{
				{Color: semiTransparentRed},
			}),
			ExtractedColors: []Color{
				{Color: red, Count: 0.4980392156862745},
			},
		},
		"Semitransparent image": {
			Image: imageFromColors([]Color{
				{Color: semiTransparentRed},
				{Color: semiTransparentRed},
				{Color: green},
			}),
			ExtractedColors: []Color{
				{Color: green, Count: 1},
				{Color: red, Count: 0.996078431372549},
			},
		},
		"Semitransparent image, bigger semitransparent region": {
			Image: imageFromColors([]Color{
				{Color: semiTransparentRed},
				{Color: semiTransparentRed},
				{Color: semiTransparentRed},
				{Color: green},
			}),
			ExtractedColors: []Color{
				{Color: red, Count: 1.4941176470588236},
				{Color: green, Count: 1},
			},
		},
		"Two colors": {
			Image: imageFromColors([]Color{
				{Color: red},
				{Color: red},
				{Color: green},
				{Color: green},
				{Color: red},
				{Color: red},
			}),
			ExtractedColors: []Color{
				{Color: red, Count: 4},
				{Color: green, Count: 2},
			},
		},
		"Mixed colors": {
			Image: imageFromColors([]Color{
				{Color: red},
				{Color: red},
				{Color: color.RGBA{R: 245, G: 0, B: 0, A: 255}},
				{Color: color.RGBA{R: 245, G: 0, B: 0, A: 255}},
				{Color: green},
				{Color: green},
				{Color: color.RGBA{R: 0, G: 240, B: 0, A: 255}},
			}),
			ExtractedColors: []Color{
				{Color: color.RGBA{R: 250, G: 0, B: 0, A: 255}, Count: 4},
				{Color: color.RGBA{R: 0, G: 250, B: 0, A: 255}, Count: 3},
			},
		},
		"File": {
			Image: imageFromFile("example/Fotolia_45549559_320_480.jpg"),
			ExtractedColors: []Color{
				{Color: color.RGBA{R: 232, G: 230, B: 228, A: 255}, Count: 30559},
				{Color: color.RGBA{R: 58, G: 58, B: 10, A: 255}, Count: 20546},
				{Color: color.RGBA{R: 205, G: 51, B: 25, A: 255}, Count: 11071},
				{Color: color.RGBA{R: 191, G: 178, B: 56, A: 255}, Count: 8705},
				{Color: color.RGBA{R: 104, G: 152, B: 12, A: 255}, Count: 5905},
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			extractedColors := ExtractColors(testCase.Image)
			if !testColorsEqual(testCase.ExtractedColors, extractedColors) {
				t.Fatalf("%v expected, got %v", testCase.ExtractedColors, extractedColors)
			}
		})
	}
}

func imageFromColors(colors []Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, len(colors), 1))
	for i, c := range colors {
		img.Set(i, 0, c.Color)
	}
	return img
}

func imageFromFile(filename string) image.Image {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}()
	img, _, _ := image.Decode(file)
	return img
}

func testColorsEqual(a, b []Color) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
