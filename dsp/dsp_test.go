package dsp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReshape(t *testing.T) {
	for i, tc := range []struct {
		input    []float64
		size     int
		expected [][]float64
	}{
		{[]float64{}, 1, [][]float64{}},
		{[]float64{1, 2, 3}, 1, [][]float64{
			[]float64{1},
			[]float64{2},
			[]float64{3},
		}},
		{[]float64{1, 2, 3}, 2, [][]float64{
			[]float64{1, 2},
		}},
		{[]float64{1, 2, 3, 4, 5, 6, 7}, 3, [][]float64{
			[]float64{1, 2, 3},
			[]float64{4, 5, 6},
		}},
		{[]float64{1, 2, 3, 4, 5, 6, 7, 8}, 4, [][]float64{
			[]float64{1, 2, 3, 4},
			[]float64{5, 6, 7, 8},
		}},
		{[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 2, [][]float64{
			[]float64{1, 2},
			[]float64{3, 4},
			[]float64{5, 6},
			[]float64{7, 8},
		}},
	} {
		require.EqualValues(t, tc.expected, Reshape(tc.input, tc.size), fmt.Sprintf("test %d", i+1))
	}
}
