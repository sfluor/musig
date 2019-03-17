package dsp

import (
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSpectrogram440(t *testing.T) {
	sampleSize := SAMPLESIZE

	s := NewSpectrogrammer(
		DOWNSAMPLERATIO,
		MAXFREQ,
		sampleSize,
	)

	file, err := os.Open("../data/440.wav")
	require.NoError(t, err)
	defer file.Close()

	spec, spr, err := s.Spectrogram(file)
	require.NoError(t, err)

	invMax := 1 / absMax(spec)

	// Threshold to remove frequencies we don't want
	threshold := 0.5
	freqBinSize := spr / DOWNSAMPLERATIO / sampleSize

	checkFreq := func(f, time, ampl float64) {
		require.InDeltaf(
			t,
			440,
			f,
			2*freqBinSize,
			"time: %f, frequency: %f, amplitude: %f",
			time,
			f,
			ampl,
		)
	}

	for time, row := range spec {
		for f := range row {
			spec[time][f] = math.Abs(spec[time][f]) * invMax
			if spec[time][f] > threshold {
				// Check that the frequency of the current point is within [440 - binSize, 440 + bin Size]
				checkFreq(float64(f)*freqBinSize, float64(time)*DOWNSAMPLERATIO/spr, spec[time][f])
			}
		}
	}

	filtered := s.HighestFreqs(spec, spr)
	for t, freqs := range filtered {
		for _, f := range freqs {
			// Use 0 as amplitude since we don't care here
			checkFreq(f, t, 0)
		}
	}
}

func TestSpectrogram440And880(t *testing.T) {
	sampleSize := SAMPLESIZE

	s := NewSpectrogrammer(
		DOWNSAMPLERATIO,
		MAXFREQ,
		sampleSize,
	)

	file, err := os.Open("../data/440_880.wav")
	require.NoError(t, err)
	defer file.Close()

	spec, spr, err := s.Spectrogram(file)
	require.NoError(t, err)

	invMax := 1 / absMax(spec)

	// Threshold to remove frequencies we don't want
	threshold := 0.5
	freqBinSize := spr / DOWNSAMPLERATIO / sampleSize

	checkFreq := func(freq, time, ampl float64) {
		// Check that the frequency of the current point is within [440 - binSize, 440 + bin Size]
		// or within [880 - binSize, 880 + bin Size]
		flag := ((freq-440) >= -freqBinSize && (freq-440) <= freqBinSize) ||
			((freq-880) >= -freqBinSize && (freq-880) <= freqBinSize)
		require.Truef(
			t,
			flag,
			"time: %f, frequency: %f, amplitude: %f",
			time,
			freq,
			ampl,
		)
	}

	for time, row := range spec {
		for f := range row {
			spec[time][f] = math.Abs(spec[time][f]) * invMax
			if spec[time][f] > threshold {
				freq := float64(f) * freqBinSize
				checkFreq(freq, float64(time)*DOWNSAMPLERATIO/spr, spec[time][f])
			}
		}
	}

	filtered := s.HighestFreqs(spec, spr)
	for t, freqs := range filtered {
		for _, f := range freqs {
			// Use 0 as amplitude since we don't care here
			checkFreq(f, t, 0)
		}
	}
}
