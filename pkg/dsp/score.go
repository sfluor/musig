package dsp

import (
	"github.com/sfluor/musig/internal/pkg/model"
	"github.com/sfluor/musig/pkg/stats"
)

// MatchScore computes a match score between the two transformed audio samples (into a list of Key + TableValue)
func MatchScore(sample, match map[model.EncodedKey]model.TableValue) float64 {
	// Will hold a list of points (time in the sample sound file, time in the matched database sound file)
	points := [2][]float64{}
	matches := 0.0
	for k, sampleValue := range sample {
		if matchValue, ok := match[k]; ok {
			points[0] = append(points[0], float64(sampleValue.AnchorTimeMs))
			points[1] = append(points[1], float64(matchValue.AnchorTimeMs))
			matches++
		}
	}
	corr := stats.Correlation(points[0], points[1])
	return corr * corr * matches
}
