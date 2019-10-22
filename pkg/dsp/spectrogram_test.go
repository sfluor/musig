package dsp

import (
	"math"
	"os"
	"path"
	"sort"
	"testing"

	"github.com/sfluor/musig/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const AssetsDir = "../../assets/test"

func TestSpectrogram440(t *testing.T) {
	sampleSize := model.SampleSize

	s := NewSpectrogrammer(
		model.DownsampleRatio,
		model.MaxFreq,
		sampleSize,
	)

	file, err := os.Open(path.Join(AssetsDir, "440.wav"))
	require.NoError(t, err)
	defer file.Close()

	spec, spr, err := s.Spectrogram(file)
	require.NoError(t, err)

	invMax := 1 / absMax(spec)

	// Threshold to remove frequencies we don't want
	threshold := 0.5
	freqBinSize := spr / model.DownsampleRatio / sampleSize

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
				checkFreq(float64(f)*freqBinSize, float64(time)*model.DownsampleRatio/spr, spec[time][f])
			}
		}
	}

	cMap := s.ConstellationMap(spec, spr)
	for _, point := range cMap {
		// Use 0 since we don't care here
		checkFreq(point.Freq, point.Time, 0)
	}

	assertIsSorted(t, cMap)
}

func TestSpectrogram440And880(t *testing.T) {
	sampleSize := model.SampleSize

	s := NewSpectrogrammer(
		model.DownsampleRatio,
		model.MaxFreq,
		sampleSize,
	)

	file, err := os.Open(path.Join(AssetsDir, "440_880.wav"))
	require.NoError(t, err)
	defer file.Close()

	spec, spr, err := s.Spectrogram(file)
	require.NoError(t, err)

	invMax := 1 / absMax(spec)

	// Threshold to remove frequencies we don't want
	threshold := 0.5
	freqBinSize := spr / model.DownsampleRatio / sampleSize

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
				checkFreq(freq, float64(time)*model.DownsampleRatio/spr, spec[time][f])
			}
		}
	}

	cMap := s.ConstellationMap(spec, spr)
	for _, point := range cMap {
		// Use 0 since we don't care here
		checkFreq(point.Freq, point.Time, 0)
	}

	assertIsSorted(t, cMap)
}

func assertIsSorted(t *testing.T, cMap []model.ConstellationPoint) {
	times := make([]float64, len(cMap))
	lastTime := 0.0
	lastFreq := 0.0
	for i, p := range cMap {
		times[i] = p.Time
		if lastTime == p.Time {
			assert.True(t, p.Freq > lastFreq)
		}
		lastTime = p.Time
		lastFreq = p.Freq
	}

	assert.True(t, sort.IsSorted(sort.Float64Slice(times)))
}
