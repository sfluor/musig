package sound

// Reader interface allows you to read audio files
type Reader interface {
	Read() ([]float64, error)
	SampleRate() uint32
}

// DOWNSAMPLERATIO is the default down sample ratio (4)
const DOWNSAMPLERATIO = 4

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
