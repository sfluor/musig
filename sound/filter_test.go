package sound

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLowPassFilter(t *testing.T) {
	// Number of samples
	N := float64(20000)
	// Sampling time
	te := float64(0.001)

	// Cutoff frequency is 10 / 2 * PI
	cutoff := 10 / (2 * math.Pi)
	lp := NewLPFilter(cutoff, 1/te)

	for _, tc := range []struct {
		f       func(float64) float64
		minAmpl float64
		maxAmpl float64
		desc    string
	}{
		// Those should not be affected
		{math.Sin, 1.98, 2, "sin(t)"},
		{math.Cos, 1.98, 2, "cos(t)"},
		{
			func(t float64) float64 { return math.Sin(t) + math.Cos(t) },
			2.8, 3, "sin(t) + cos(t)",
		},
		// Those should be filtered out
		{
			func(t float64) float64 { return math.Sin(t * 1000) },
			0, 0.05, "sin(1000 * t)",
		},
		{
			func(t float64) float64 { return math.Cos(t * 1000) },
			0, 0.05, "cos(1000 * t)",
		},
		// Those should be splitted in two
		{
			func(t float64) float64 { return math.Sin(t) + 10*math.Sin(t*1000) },
			1.98, 2.2, "sin(t) + 50 * sin(t * 10e5)",
		},
		{
			func(t float64) float64 { return math.Cos(t) + 10*math.Cos(t*1000) },
			1.98, 2.2, "cos(t) + 50 * cos(t * 10e5)",
		},
		// This should not be affected
		{
			func(t float64) float64 { return 1 },
			0.99, 1, "constant",
		},
	} {
		values := make([]float64, 0, int(N))

		for n := float64(0); n < N; n++ {
			values = append(values, tc.f(n*te))
		}
		ampl := sigAmpl(lp.Filter(values))
		require.Truef(t, ampl <= tc.maxAmpl, "max amplitude: expected %f (filtered) < %f (max) for '%s'", ampl, tc.maxAmpl, tc.desc)
		require.Truef(t, tc.minAmpl <= ampl, "min amplitude: expected %f (min) < %f (filtered) for '%s'", tc.minAmpl, ampl, tc.desc)
	}
}

func sigAmpl(arr []float64) float64 {
	max, min := arr[0], arr[0]
	for _, x := range arr {
		max = math.Max(max, x)
		min = math.Min(min, x)
	}
	return max - min
}
