package sound

import (
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
