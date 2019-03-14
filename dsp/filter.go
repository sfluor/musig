package dsp

import (
	"math"
)

// MAXFREQ is 5kHz
const MAXFREQ = float64(5000)

// Filterer is the interface used for filters
type Filterer interface {
	Filter(arr []float64) []float64
}

// LPFilter is a first order low pass filter usig H(p) = 1 / (1 + pRC)
type LPFilter struct {
	alpha float64
	rc    float64
}

// NewLPFilter creates a new low pass Filter
func NewLPFilter(cutoff, sampleRate float64) *LPFilter {
	rc := 1 / (cutoff * 2 * math.Pi)
	dt := 1 / sampleRate
	alpha := dt / (rc + dt)
	return &LPFilter{alpha: alpha, rc: rc}
}

// XXX: This filter could be improved
// Filter filters the given array of values, panics if array is empty
func (lp *LPFilter) Filter(arr []float64) []float64 {
	res := make([]float64, 0, len(arr))

	// First value should be arr[0] / RC using the initial value theorem
	res = append(res, arr[0]*lp.alpha)

	for i := range arr[:len(arr)-1] {
		res = append(
			res,
			// Formula used is:
			// y(n+1) = y(n) + alpha * (x(n+1) -y(n))
			res[i]+lp.alpha*(arr[i+1]-res[i]),
		)
	}
	return res
}
