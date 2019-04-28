package dsp

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFFT(t *testing.T) {
	for _, fft := range []func([]float64) []float64{
		FFT, DFT,
	} {
		for j, tc := range []struct {
			input    []float64
			expected []float64
		}{
			{
				[]float64{1, 0, 0, 0, 0, 0, 0, 0},
				[]float64{1, 1, 1, 1, 1, 1, 1, 1},
			},
			{
				[]float64{1, 2},
				[]float64{3, -1},
			},
			{
				[]float64{1, 2, 3},
				[]float64{6, -1.5, -1.5},
			},
			{
				[]float64{-1, 0, 1, 0},
				[]float64{0, -2, 0, -2},
			},
			{
				[]float64{1, 2, 3, 4, 5, 6},
				[]float64{21, -3, -3, -3, -3, -3},
			},
		} {
			output := fft(tc.input)
			require.Equal(t, len(tc.expected), len(output), fmt.Sprintf("test %d", j+1))
			require.InDeltaSlice(t, tc.expected, output, 0.01, fmt.Sprintf("test %d", j+1))
		}
	}
}

func BenchmarkFT(b *testing.B) {
	times := []int{10, 50, 100, 500, 1000, 5000}
	inputs := map[int][]float64{}
	f := func(f float64) float64 { return math.Sin(2*f) + math.Cos(3*f) }

	// Compute inputs
	for _, N := range times {
		inputs[N] = make([]float64, 0, N)
		for i := 0; i < N; i++ {
			inputs[N] = append(inputs[N], f(float64(i)))
		}
	}

	b.ResetTimer()
	for N, input := range inputs {
		b.Run(fmt.Sprintf("DFT_%d", N), benchOneFT(b, input, DFT))
		b.Run(fmt.Sprintf("FFT_%d", N), benchOneFT(b, input, FFT))
	}
}

func benchOneFT(b *testing.B, input []float64, fft func([]float64) []float64) func(b *testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fft(input)
		}
	}
}
