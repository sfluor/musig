package dsp

import "math"

// WindowFunc represents a window function
// It takes the number of points we want in the output window
type WindowFunc func(M int) []float64

// ApplyWindow applies the provided window on the given array
func ApplyWindow(arr []float64, wf WindowFunc) {
	for i, v := range wf(len(arr)) {
		arr[i] *= v
	}
}

// HammingWindow is a hamming window
// the formula used is w(n) = 0.54 - 0.46 * cos(2 * pi * n / (M - 1)) for M in [0, M-1]
func HammingWindow(M int) []float64 {
	switch M {
	case 0:
		return []float64{}
	case 1:
		return []float64{1}
	default:
		f := 2 * math.Pi / float64(M-1)
		res := make([]float64, M)
		for n := range res {
			res[n] = 0.54 - 0.46*math.Cos(f*float64(n))
		}
		return res
	}
}
