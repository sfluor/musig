package stats

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAbsMax(t *testing.T) {
	for _, tc := range []struct {
		input []float64
		max   float64
		idx   int
	}{
		{
			[]float64{-1.5, 1, -0.2, 1.45},
			1.5,
			0,
		},
		{
			[]float64{-100, 101},
			101,
			1,
		},
		{
			[]float64{-100, 100, 50},
			100,
			0,
		},
	} {
		m, i := ArgAbsMax(tc.input)
		require.Equal(t, tc.max, m)
		require.Equal(t, tc.idx, i)
	}
}

func TestAvg(t *testing.T) {
	for _, tc := range []struct {
		input    []float64
		expected float64
	}{
		{
			[]float64{-1.5, 1, -0.2, 1.45},
			0.1875,
		},
		{
			[]float64{-100, 101},
			0.5,
		},
		{
			[]float64{-100, 100, 50},
			16.6666,
		},
	} {
		require.InDelta(t, tc.expected, Avg(tc.input), 0.01)
	}
}

func TestArgAbove(t *testing.T) {
	for _, tc := range []struct {
		input     []float64
		threshold float64
		expected  []int
	}{
		{
			[]float64{-1.5, 1, -0.2, 1.45},
			-2,
			[]int{0, 1, 2, 3},
		},
		{
			[]float64{-1.5, 1, -0.2, 1.45},
			0,
			[]int{1, 3},
		},
		{
			[]float64{-100, 101},
			101.1,
			nil,
		},
		{
			[]float64{-100, 100, 50},
			25,
			[]int{1, 2},
		},
	} {
		require.EqualValues(t, tc.expected, ArgAbove(tc.threshold, tc.input))
	}
}
