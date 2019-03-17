package sound

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	riff "github.com/youpy/go-riff"
	wav "github.com/youpy/go-wav"
)

var _ Reader = &WAVReader{}

// WAVReader implements the sound.Reader interface
type WAVReader struct {
	wr         *wav.Reader
	isStereo   bool
	sampleRate float64
}

// NewWAVReader creates a new wav reader
func NewWAVReader(r riff.RIFFReader) (*WAVReader, error) {
	reader := wav.NewReader(r)
	f, err := reader.Format()
	if err != nil {
		return nil, errors.Wrap(err, "could not get wav format")
	}

	return &WAVReader{
		wr:         reader,
		isStereo:   f.NumChannels == 2,
		sampleRate: float64(f.SampleRate),
	}, nil
}

// Read reads from the given wav file and return raw audio data or an error if an error occured
func (r *WAVReader) Read(dst []float64, N int) (int, error) {
	if len(dst) != N {
		return 0, fmt.Errorf("given dst has size %d, expected %d", len(dst), N)
	}

	// go-wav uses 4 * the number of samples we want to read as parameter
	samples, err := r.wr.ReadSamples(uint32(N))
	if err == io.EOF {
		return 0, err
	}
	if err != nil {
		return 0, errors.Wrap(err, "could not read samples")
	}

	// Take care of mono / stereo
	// If sound is in stereo we want
	size := len(samples)
	if r.isStereo {
		size *= 2
	}

	if r.isStereo {
		for i, sample := range samples {
			// We average the two entries in case of stereo
			dst[i] = float64(sample.Values[0]+sample.Values[1]) / 2
		}
		return size / 2, nil
	}

	for i, sample := range samples {
		dst[i] = float64(sample.Values[0])
	}

	return size, nil
}

// SampleRate returns the sample rate for the given reader
func (r *WAVReader) SampleRate() float64 {
	return r.sampleRate
}
