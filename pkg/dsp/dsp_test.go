package dsp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownsample(t *testing.T) {
	for _, tc := range []struct {
		input    []float64
		expected []float64
		ratio    int
		desc     string
	}{
		{
			[]float64{1, 2, 3, 4, 5},
			[]float64{1, 2, 3, 4, 5},
			1,
			"1-ratio (no downsampling)",
		},
		{
			[]float64{1, 2, 3, 4, 5},
			[]float64{1.5, 3.5, 5},
			2,
			"2-ratio",
		},
		{
			[]float64{1, 2, 3, 4, 5},
			[]float64{2, 4.5},
			3,
			"3-ratio",
		},
		{
			[]float64{1, 2, 3, 4, 5},
			[]float64{2.5, 5},
			4,
			"4-ratio",
		},
		{
			[]float64{1, 2, 3, 4, 5},
			[]float64{3},
			5,
			"5-ratio",
		},
	} {
		output := Downsample(tc.input, tc.ratio)
		require.EqualValues(t, tc.expected, output, tc.desc)
	}
}

func TestReshape(t *testing.T) {
	for i, tc := range []struct {
		input    []float64
		size     int
		expected [][]float64
	}{
		{[]float64{}, 1, [][]float64{}},
		{[]float64{1, 2, 3}, 1, [][]float64{
			{1},
			{2},
			{3},
		}},
		{[]float64{1, 2, 3}, 2, [][]float64{
			{1, 2},
		}},
		{[]float64{1, 2, 3, 4, 5, 6, 7}, 3, [][]float64{
			{1, 2, 3},
			{4, 5, 6},
		}},
		{[]float64{1, 2, 3, 4, 5, 6, 7, 8}, 4, [][]float64{
			{1, 2, 3, 4},
			{5, 6, 7, 8},
		}},
		{[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 2, [][]float64{
			{1, 2},
			{3, 4},
			{5, 6},
			{7, 8},
		}},
	} {
		require.EqualValues(t, tc.expected, Reshape(tc.input, tc.size), fmt.Sprintf("test %d", i+1))
	}
}
