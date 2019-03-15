package dsp

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFFT(t *testing.T) {
	for _, fft := range []func([]complex128) []complex128{
		FFT, DFT,
	} {
		for j, tc := range []struct {
			input    []complex128
			expected []complex128
		}{
			{
				[]complex128{1, 0, 0, 0, 0, 0, 0, 0},
				[]complex128{1, 1, 1, 1, 1, 1, 1, 1},
			},
			{
				[]complex128{1, 2},
				[]complex128{3, -1},
			},
			{
				[]complex128{1, 2, 3},
				[]complex128{6, -1.5 + 0.866i, -1.5 - 0.866i},
			},
			{
				[]complex128{-1, 0, 1, 0},
				[]complex128{0, -2, 0, -2},
			},
			{
				[]complex128{1, 2, 3, 4, 5, 6},
				[]complex128{21, -3 + 5.196i, -3 + 1.732i, -3, -3 - 1.732i, -3 - 5.196i},
			},
		} {
			output := fft(tc.input)
			require.Equal(t, len(tc.expected), len(output), fmt.Sprintf("test %d", j+1))
			for i, e := range tc.expected {
				require.InDelta(t, real(e), real(output[i]), 0.01, fmt.Sprintf("test %d, index: %d", j+1, i+1))
				require.InDelta(t, imag(e), imag(output[i]), 0.01, fmt.Sprintf("test %d, index: %d", j+1, i+1))
			}
		}
	}
}

func BenchmarkFT(b *testing.B) {
	times := []int{10, 50, 100, 500, 1000, 5000}
	inputs := map[int][]complex128{}
	f := func(f float64) complex128 { return complex(math.Sin(2*f)+math.Cos(3*f), 0) }

	// Compute inputs
	for _, N := range times {
		inputs[N] = make([]complex128, 0, N)
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

func benchOneFT(b *testing.B, input []complex128, fft func([]complex128) []complex128) func(b *testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fft(input)
		}
	}
}
