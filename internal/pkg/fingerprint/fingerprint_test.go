package fingerprint

import (
	"os"
	"path"
	"testing"

	"github.com/sfluor/musig/internal/pkg/model"
	"github.com/sfluor/musig/pkg/dsp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const AssetsDir = "../../../assets/test"

func TestFingerprinting440(t *testing.T) {
	testFingerprintingOnFile(t, path.Join(AssetsDir, "440.wav"))
}

func TestFingerprinting440And880(t *testing.T) {
	testFingerprintingOnFile(t, path.Join(AssetsDir, "440_880.wav"))
}

func testFingerprintingOnFile(t *testing.T, path string) {
	sampleSize := model.SAMPLESIZE

	s := dsp.NewSpectrogrammer(
		model.DOWNSAMPLERATIO,
		model.MAXFREQ,
		sampleSize,
	)

	file, err := os.Open(path)
	require.NoError(t, err)
	defer file.Close()

	spec, spr, err := s.Spectrogram(file)
	require.NoError(t, err)

	f := NewDefaultFingerprinter()

	cMap := s.ConstellationMap(spec, spr)

	// Apply a second constellation map only on a sub part of the file
	subSpec := spec[40:60]
	subCMap := s.ConstellationMap(subSpec, spr)

	table := f.Fingerprint(0, cMap)
	subTable := f.Fingerprint(0, subCMap)

	for key := range subTable {
		assert.Contains(t, table, key)
	}
}
