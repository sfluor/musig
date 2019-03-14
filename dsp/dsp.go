package dsp

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
