package dsp

import (
	"math"
	"math/cmplx"

	"gonum.org/v1/gonum/fourier"
)

// DFT is a discrete fourier transform implementation (it is slow O(N^2))
// TODO parallelize ?
func DFT(arr []complex128) []complex128 {
	N := len(arr)
	res := make([]complex128, N)
	theta := -2i * math.Pi / itoc(N)
	for n := range res {
		for k := 0; k < N; k++ {
			res[n] += arr[k] * cmplx.Exp(theta*itoc(k)*itoc(n))
		}
	}
	return res
}

// FFT is a fast fourier transform using gonum/fourier
// TODO remove adding duplicate frequencies (it's ~33% slower with them)
func FFT(arr []complex128) []complex128 {
	fft := fourier.NewFFT(len(arr))
	in := make([]float64, len(arr))

	for i, v := range arr {
		in[i] = real(v)
	}

	res := fft.Coefficients(nil, in)

	// Add duplicate frequencies
	for i := len(arr) - len(res); i >= 1; i-- {
		res = append(res, cmplx.Conj(res[i]))
	}

	return res
}

// Int to complex
func itoc(i int) complex128 {
	return complex(float64(i), 0)
}
