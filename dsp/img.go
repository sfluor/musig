package dsp

import (
	"image"
	"image/color"
	"math"

	"github.com/sfluor/musig/stats"
)

// SpecToImg takes a spectrogram (matrix of floats representing m[time][frequency] = amplitude)
// and return an image
func SpecToImg(matrix [][]float64) image.Image {
	// TODO stop hardcoding this
	// 10 Pixel for 1 frequency step
	// 150 Pixels for 1 timeStep
	nTime, nFreq := 150, 20

	img := image.NewRGBA(image.Rect(0, 0, nTime*len(matrix), nFreq*len(matrix[0])))

	// XXX the way values are normalized could be changed
	invMax := 1 / absMax(matrix)
	for time, row := range matrix {
		for freq, a := range row {
			color := colorbar(math.Abs(a) * invMax)
			for t := 0; t < nTime; t++ {
				for f := 0; f < nFreq; f++ {
					img.Set(
						// x
						nTime*time+t,
						// y
						nFreq*len(matrix[0])-nFreq*freq+f,
						// color
						color,
					)
				}
			}
		}
	}

	return img
}

// absMax returns the maximum value of abs(mat) where mat is the given matrix
func absMax(mat [][]float64) (max float64) {
	for _, row := range mat {
		max = math.Max(max, stats.AbsMax(row))
	}
	return max
}

// colorbar is function to map a value in [0, 1] to a color
func colorbar(val float64) color.RGBA {
	r := 255 * math.Min(math.Max(0, 1.5-math.Abs(1-4*(val-0.5))), 1)
	g := 255 * math.Min(math.Max(0, 1.5-math.Abs(1-4*(val-0.25))), 1)
	b := 255 * math.Min(math.Max(0, 1.5-math.Abs(1-4*val)), 1)
	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
}
