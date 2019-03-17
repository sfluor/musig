package dsp

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/sfluor/musig/sound"
	"github.com/sfluor/musig/stats"
)

// Spectrogrammer is a struct that allows to create spectrograms
type Spectrogrammer struct {
	dsRatio float64
	maxFreq float64
	binSize float64
}

// NewSpectrogrammer creates a new spectrogrammer
func NewSpectrogrammer(dsRatio, maxFreq, binSize float64) *Spectrogrammer {
	return &Spectrogrammer{
		dsRatio: dsRatio,
		maxFreq: maxFreq,
		binSize: binSize,
	}
}

// Spectrogram reads the provided audio file and returns a spectrogram for it
// Matrix is in the following format:
// TIME : FREQUENCY : Value
// time is t * dsp.DOWNSAMPLERATIO / reader.SampleRate()
// frequency is f * freqBinSize
func (s *Spectrogrammer) Spectrogram(file *os.File) ([][]float64, float64, error) {
	reader, err := sound.NewWAVReader(file)
	if err != nil {
		return nil, 0, errors.Wrap(err, "error reading wav")
	}

	spr := reader.SampleRate()
	lp := NewLPFilter(s.maxFreq, spr)
	matrix := [][]float64{}

	bin := make([]float64, int(s.binSize*s.dsRatio))
	for {
		n, err := reader.Read(bin, int(s.binSize*s.dsRatio))
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, 0, errors.Wrap(err, "error reading from sound file")
		}

		// TODO handle this edge case
		if n != int(s.binSize*s.dsRatio) {
			break
		}

		fft := FFT(
			Downsample(
				lp.Filter(bin[:n]),
				int(s.dsRatio),
			),
		)

		// TODO remove slicing here when removing duplicate values retuned by the FFT
		matrix = append(matrix, fft[:len(fft)-(len(fft)-1)/2-1])
	}

	s.HighestFreqs(matrix, reader.SampleRate())
	return matrix, spr, nil
}

// HighestFreqs takes a spectrogram, it's sample rate and returns the highest frequencies and their time in the audio file
func (s *Spectrogrammer) HighestFreqs(spec [][]float64, sampleRate float64) map[float64][]float64 {
	// TODO stop hardcoding those
	coef := 2.0
	// For each 512-sized bins create logarithmic bands
	// [0, 10], [10, 20], [20, 40], [40, 80], [80, 160], [160, 511]
	bands := [][]int{{0, 10}, {10, 20}, {20, 40}, {40, 80}, {80, 160}, {160, 512}}

	res := map[float64][]float64{}

	// Frequency bin size
	fbs := s.freqBinSize(sampleRate)

	// Maximum of amplitude and their corresponding frequencies
	var maxs, freqs []float64
	var idx int
	for t, row := range spec {
		maxs, freqs = make([]float64, len(bands)), make([]float64, len(bands))

		// We retrieve the maximum amplitudes and their frequency
		for i, band := range bands {
			maxs[i], idx = stats.ArgAbsMax(row[band[0]:band[1]])
			freqs[i] = float64(band[0]+idx) * fbs
		}

		// Keep only the bins above the average of the max bins
		avg := stats.Avg(maxs)
		indices := stats.ArgAbove(avg*coef, maxs)

		// Register the frequencies we kept and their time of apparition
		time := float64(t) / sampleRate
		res[time] = make([]float64, len(indices))

		for i, idx := range indices {
			res[time][i] = freqs[idx]
		}
	}

	return res
}

// Returns the bin size for a frequency bin given a sample rate
func (s *Spectrogrammer) freqBinSize(spr float64) float64 {
	return spr / s.dsRatio / s.binSize
}
