package sound

import (
	"io"

	"github.com/gordonklaus/portaudio"
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
func (r *WAVReader) Read(dst []float64) (int, error) {
	N := uint32(len(dst))

	// go-wav uses 4 * the number of samples we want to read as parameter
	samples, err := r.wr.ReadSamples(N)
	if err == io.EOF {
		return 0, err
	}
	if err != nil {
		return 0, errors.Wrap(err, "could not read samples")
	}

	// Take care of mono / stereo
	// If sound is in stereo we want to get it into mono
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

// RecordWAV listens on the microphone and saves the signal to the given io.Writer
// it takes a stop channel to interrupt the recording
func RecordWAV(writer io.Writer, stopCh <-chan struct{}) error {
	samples := []wav.Sample{}
	channels := 1
	bitsPerSample := 32
	sampleRate := 44100

	portaudio.Initialize()
	defer portaudio.Terminate()

	in := make([]int32, 64)
	stream, err := portaudio.OpenDefaultStream(channels, 0, float64(sampleRate), len(in), in)
	if err != nil {
		return errors.Wrap(err, "error opening portaudio stream")
	}
	defer stream.Close()

	err = stream.Start()
	if err != nil {
		return errors.Wrap(err, "error starting stream")
	}

listen:
	for {
		err = stream.Read()
		if err != nil {
			return errors.Wrap(err, "error reading from stream")
		}

		for _, v := range in {
			samples = append(samples, wav.Sample{Values: [2]int{
				// Append the same value twice, worst case it will be used as stereo and averaged
				int(v), int(v),
			}})
		}

		select {
		case <-stopCh:
			break listen
		default:
		}
	}

	if err = stream.Stop(); err != nil {
		return errors.Wrap(err, "error stoping stream")
	}

	wr := wav.NewWriter(writer, uint32(len(samples)), uint16(channels), uint32(sampleRate), uint16(bitsPerSample))
	if err = wr.WriteSamples(samples); err != nil {
		return errors.Wrap(err, "error writing samples")
	}

	return nil
}
