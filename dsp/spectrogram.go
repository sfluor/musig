package dsp

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/sfluor/musig/sound"
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

		fft := FFT(
			Downsample(
				lp.Filter(bin[:n]),
				int(s.dsRatio),
			),
		)

		// TODO remove slicing here when removing duplicate values retuned by the FFT
		matrix = append(matrix, fft[:len(fft)-(len(fft)-1)/2])
	}

	return matrix, spr, nil
}
