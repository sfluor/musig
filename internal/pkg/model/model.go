package model

import (
	"fmt"
	"path/filepath"
	"strings"
)

// DownsampleRatio is the default down sample ratio (4)
const DownsampleRatio = 4

// SampleSize is the default sample size (1024)
const SampleSize = 1024.0

// MaxFreq is 5kHz
const MaxFreq = 5000.0

// ConstellationPoint represents a point in the constellation map (time + frequency)
type ConstellationPoint struct {
	Time float64
	Freq float64
}

func (cp ConstellationPoint) String() string {
	return fmt.Sprintf("(t: %f, f: %f)", cp.Time, cp.Freq)
}

// SongNameFromPath returns the song file name from the given path (removing the path and the extension)
func SongNameFromPath(path string) string {
	ext := filepath.Ext(path)
	return strings.TrimRight(filepath.Base(path), ext)
}
