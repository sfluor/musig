package dsp

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/sfluor/musig/internal/pkg/model"
	"github.com/sfluor/musig/internal/pkg/sound"
	"github.com/sfluor/musig/pkg/stats"
)

// Spectrogrammer is a struct that allows to create spectrograms
type Spectrogrammer struct {
	dsRatio float64
	maxFreq float64
	binSize float64
	// thresholdCoefficient is used to filter out the important frequencies
	// increasing it decreases the size of the constellation maps returned
	thresholdCoefficient float64
	// windowing is used to activate / deactivate the windowing function
	windowing bool
}

// NewSpectrogrammer creates a new spectrogrammer
func NewSpectrogrammer(dsRatio, maxFreq, binSize float64, windowing bool) *Spectrogrammer {
	return &Spectrogrammer{
		dsRatio: dsRatio,
		maxFreq: maxFreq,
		binSize: binSize,
		// TODO stop hardcoding this
		thresholdCoefficient: 1,
		windowing: windowing,
	}
}

// Spectrogram reads the provided audio file and returns a spectrogram for it
// Matrix is in the following format:
// TIME : FREQUENCY : Value
// time is t * binSize * dsp.DownsampleRatio / reader.SampleRate()
// frequency is f * freqBinSize
func (s *Spectrogrammer) Spectrogram(file *os.File) ([][]float64, float64, error) {
	reader, err := sound.NewWAVReader(file)
	if err != nil {
		return nil, 0, errors.Wrap(err, "error reading wav")
	}

	spr := reader.SampleRate()
	lp := NewLPFilter(s.maxFreq, spr)
	var matrix [][]float64

	bin := make([]float64, int(s.binSize*s.dsRatio))
	for {
		n, err := reader.Read(bin)
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

		sampled := Downsample(
			lp.Filter(bin[:n]),
			int(s.dsRatio),
		)
		if s.windowing {
			ApplyWindow(sampled, HammingWindow)
		}
		fft := FFT(sampled)

		// TODO remove slicing here when removing duplicate values returned by the FFT
		matrix = append(matrix, fft[:len(fft)-(len(fft)-1)/2-1])
	}

	return matrix, spr, nil
}

// ConstellationMap takes a spectrogram, its sample rate and returns the highest frequencies and their time in the audio file
// The returned slice is ordered by time and is ordered by frequency for a constant time:
// If two time-frequency points have the same time, the time-frequency point with the lowest frequency is before the other one.
// If a time time-frequency point has a lower time than another point one then it is before.
func (s *Spectrogrammer) ConstellationMap(spec [][]float64, sampleRate float64) []model.ConstellationPoint {
	// For each 512-sized bins create logarithmic bands
	// [0, 10], [10, 20], [20, 40], [40, 80], [80, 160], [160, 511]
	bands := [][]int{{0, 10}, {10, 20}, {20, 40}, {40, 80}, {80, 160}, {160, 512}}

	var res []model.ConstellationPoint

	// Frequency bin size
	fbs := s.freqBinSize(sampleRate)

	// Time step
	timeStep := s.dsRatio * s.binSize / sampleRate

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

		// Keep only the bins above the average times the threshold coeff of the max bins
		avg := stats.Avg(maxs)
		indices := stats.ArgAbove(avg*s.thresholdCoefficient, maxs)

		// Register the frequencies we kept and their time of apparition
		time := timeStep * float64(t)

		for _, idx := range indices {
			res = append(res, model.ConstellationPoint{Time: time, Freq: freqs[idx]})
		}
	}

	return res
}

// Returns the bin size for a frequency bin given a sample rate
func (s *Spectrogrammer) freqBinSize(spr float64) float64 {
	return spr / s.dsRatio / s.binSize
}
