package utils

import (
	"fmt"
	"image"
	"math"

	"github.com/disintegration/imaging"
)

// img.At(x, y).RGBA() returns four uint32 values, we want int
func rgbaToInt(r uint32, g uint32, b uint32, a uint32) (int, int, int, int) {
	return int(r / 257), int(g / 257), int(b / 257), int(a / 257)
}

func imageSimilarityIndexFile(origin, reference string, sensitivity float32) (float32, error) {
	img, err := imaging.Open(origin)
	if err != nil {
		return 0, fmt.Errorf("Error opening image file: %s", origin)
	}
	img = imaging.Blur(img, 3.5)
	ref, err := imaging.Open(reference)
	if err != nil {
		return 0, fmt.Errorf("Error opening reference image file: %s", reference)
	}
	ref = imaging.Blur(ref, 3.5)
	return imageSimilarityIndex(img, ref, sensitivity)
}

// Difference produces a diff beteen two images.
func imageSimilarityIndex(img, ref image.Image, sensitivity float32) (float32, error) {
	if img.Bounds() != ref.Bounds() {
		return 0, fmt.Errorf("Images size do not match")
	}

	bounds := img.Bounds()
	da := make([]float64, bounds.Max.Y)

	parallel(0, bounds.Max.Y, func(ys <-chan int) {
		for y := range ys {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				ir, ig, ib, _ := rgbaToInt(img.At(x, y).RGBA())
				rr, rg, rb, _ := rgbaToInt(ref.At(x, y).RGBA())
				f := 0.299*float64(ir) + 0.587*float64(ig) + 0.114*float64(ib)
				g := 0.299*float64(rr) + 0.587*float64(rg) + 0.114*float64(rb)
				da[y] += (f - g) * (f - g) / 65535.0
			}
		}
	})

	var sum float64
	for _, v := range da {
		sum += v
	}

	sum = math.Pow(sum/float64(bounds.Max.X*bounds.Max.Y), float64(sensitivity))

	return float32(sum), nil
}
