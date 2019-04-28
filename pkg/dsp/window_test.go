package dsp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHammingWindow(t *testing.T) {
	for _, tc := range []struct {
		m        int
		expected []float64
	}{
		{0, []float64{}},
		{1, []float64{1}},
		{2, []float64{0.08, 0.08}},
		{5, []float64{0.08, 0.54, 1., 0.54, 0.08}},
		{11, []float64{0.08, 0.16785218, 0.39785218, 0.68214782, 0.91214782, 1., 0.91214782, 0.68214782, 0.39785218, 0.16785218, 0.08}},
	} {
		require.InDeltaSlice(t, tc.expected, HammingWindow(tc.m), 0.0001)
	}
}

func TestApplyHammingWindow(t *testing.T) {
	for i, tc := range []struct {
		input    []float64
		expected []float64
	}{
		{[]float64{}, []float64{}},
		{[]float64{1}, []float64{1}},
		{[]float64{1, 1}, []float64{0.08, 0.08}},
		{[]float64{10, 2, 0, 5, 4}, []float64{0.8, 1.08, 0, 2.7, 0.32}},
		{[]float64{0, 0, 0}, []float64{0, 0, 0}},
	} {
		ApplyWindow(tc.input, HammingWindow)
		require.InDeltaSlice(t, tc.expected, tc.input, 0.0001, fmt.Sprintf("test %d", i+1))
	}
}
