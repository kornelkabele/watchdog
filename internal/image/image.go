package utils

import (
	"fmt"
	"image"
	"math"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/disintegration/imaging"
)

var maxProcs int64

// SetMaxProcs limits the number of concurrent processing goroutines to the given value.
// A value <= 0 clears the limit.
func SetMaxProcs(value int) {
	atomic.StoreInt64(&maxProcs, int64(value))
}

// parallel processes the data in separate goroutines.
func parallel(start, stop int, fn func(<-chan int)) {
	count := stop - start
	if count < 1 {
		return
	}

	procs := runtime.GOMAXPROCS(0)
	limit := int(atomic.LoadInt64(&maxProcs))
	if procs > limit && limit > 0 {
		procs = limit
	}
	if procs > count {
		procs = count
	}

	c := make(chan int, count)
	for i := start; i < stop; i++ {
		c <- i
	}
	close(c)

	var wg sync.WaitGroup
	for i := 0; i < procs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn(c)
		}()
	}
	wg.Wait()
}

// img.At(x, y).RGBA() returns four uint32 values, we want int
func rgbaToInt(r uint32, g uint32, b uint32, a uint32) (int, int, int, int) {
	return int(r / 257), int(g / 257), int(b / 257), int(a / 257)
}

// ImageSimilarityIndexFile produces a diff beteen two image files.
func ImageSimilarityIndexFile(origin, reference string, sensitivity float32) (float32, error) {
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
	return ImageSimilarityIndex(img, ref, sensitivity)
}

// ImageSimilarityIndex produces a diff beteen two images.
func ImageSimilarityIndex(img, ref image.Image, sensitivity float32) (float32, error) {
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
