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
	semiTransparentRedTransformed := color.RGBA{R: 127, G: 0, B: 0, A: 255}

	testCases := map[string]struct {
		Image           image.Image
		ExtractedColors []Color
	}{
		"Empty file": {
			Image:           imageFromColors([]color.Color{}),
			ExtractedColors: []Color{},
		},
		"Single pixel": {
			Image: imageFromColors([]color.Color{
				red,
			}),
			ExtractedColors: []Color{
				{Color: red, Count: 1},
			},
		},
		"One color": {
			Image: imageFromColors([]color.Color{
				white,
				white,
				white,
				white,
				white,
			}),
			ExtractedColors: []Color{
				{Color: white, Count: 5},
			},
		},
		"Transparent image": {
			Image: imageFromColors([]color.Color{
				white,
				white,
				white,
				transparent,
			}),
			ExtractedColors: []Color{
				{Color: white, Count: 3},
			},
		},
		"Semitransparent single pixel": {
			Image: imageFromColors([]color.Color{
				semiTransparentRed,
			}),
			ExtractedColors: []Color{
				{Color: semiTransparentRedTransformed, Count: 1},
			},
		},
		"Semitransparent image": {
			Image: imageFromColors([]color.Color{
				semiTransparentRed,
				semiTransparentRed,
				green,
			}),
			ExtractedColors: []Color{
				{Color: semiTransparentRedTransformed, Count: 2},
				{Color: green, Count: 1},
			},
		},
		"Semitransparent image, bigger semitransparent region": {
			Image: imageFromColors([]color.Color{
				semiTransparentRed,
				semiTransparentRed,
				semiTransparentRed,
				green,
			}),
			ExtractedColors: []Color{
				{Color: semiTransparentRedTransformed, Count: 3},
				{Color: green, Count: 1},
			},
		},
		"Two colors": {
			Image: imageFromColors([]color.Color{
				red,
				red,
				green,
				green,
				red,
				red,
			}),
			ExtractedColors: []Color{
				{Color: red, Count: 4},
				{Color: green, Count: 2},
			},
		},
		"Mixed colors": {
			Image: imageFromColors([]color.Color{
				red,
				red,
				color.RGBA{R: 245, G: 0, B: 0, A: 255},
				color.RGBA{R: 245, G: 0, B: 0, A: 255},
				green,
				green,
				color.RGBA{R: 0, G: 240, B: 0, A: 255},
			}),
			ExtractedColors: []Color{
				{Color: color.RGBA{R: 245, G: 0, B: 0, A: 255}, Count: 4},
				{Color: color.RGBA{R: 0, G: 240, B: 0, A: 255}, Count: 3},
			},
		},
		"File": {
			Image: imageFromFile("example/Fotolia_45549559_320_480.jpg"),
			ExtractedColors: []Color{
				{Color: color.RGBA{R: 213, G: 196, B: 154, A: 255}, Count: 30559},
				{Color: color.RGBA{R: 75, G: 92, B: 2, A: 255}, Count: 20546},
				{Color: color.RGBA{R: 135, G: 0, B: 0, A: 255}, Count: 11071},
				{Color: color.RGBA{R: 131, G: 131, B: 57, A: 255}, Count: 8705},
				{Color: color.RGBA{R: 124, G: 164, B: 3, A: 255}, Count: 5905},
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

func imageFromColors(colors []color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, len(colors), 1))
	for i, c := range colors {
		img.Set(i, 0, c)
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
