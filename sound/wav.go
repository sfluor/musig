package sound

import (
	"io"

	"github.com/pkg/errors"
	riff "github.com/youpy/go-riff"
	wav "github.com/youpy/go-wav"
)

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
func (r *WAVReader) Read() ([]float64, error) {
	samples, err := r.wr.ReadSamples()
	if err == io.EOF {
		return nil, err
	}
	if err != nil {
		return nil, errors.Wrap(err, "could not read samples")
	}

	// Take care of mono / stereo
	// If sound is in stereo we want
	size := len(samples)
	if r.isStereo {
		size *= 2
	}

	res := make([]float64, 0, size)
	for _, sample := range samples {
		// TODO: This is not efficient, split this ?
		if r.isStereo {
			// We average the two entries in case of stereo
			res = append(res, (float64(sample.Values[0]+sample.Values[1]))/2)
		} else {
			res = append(res, float64(sample.Values[0]))
		}
	}

	return res, nil
}

// SampleRate returns the sample rate for the given reader
func (r *WAVReader) SampleRate() float64 {
	return r.sampleRate
}
