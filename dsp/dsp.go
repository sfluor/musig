package dsp

// Downsample downsamples the given array using the given ratio by averaging
func Downsample(arr []float64, ratio int) []float64 {
	res := make([]float64, 0, len(arr)/ratio)

	// Used to divide the samples (might be different than ratio if len(arr) % ratio != 0
	div := float64(0)

	for i := 0; i < len(arr); i += ratio {
		res = append(res, 0)
		for j := 0; j < ratio && i+j < len(arr); j++ {
			div++
			res[i/ratio] += arr[i+j]
		}
		res[i/ratio] /= div
		div = 0
	}

	return res
}

// Reshape takes an array and a bin size and will split the given array
// into multiple bins of the given bin size
// if len(arr) % size != 0, the data at the end will be dropped
// XXX: maybe don't drop it ?
func Reshape(arr []float64, size int) [][]float64 {
	res := make([][]float64, len(arr)/size)
	for i := range res {
		res[i] = arr[i*size : (i+1)*size]
	}
	return res
}
