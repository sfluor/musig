package stats

import (
	"math"
)

// Correlation computes the correlation between 2 series of points
// the length used is x's
func Correlation(x []float64, y []float64) float64 {
	n := len(x)
	meanX, meanY := Avg(x[:n]), Avg(y[:n])

	sXY := 0.0
	sX := 0.0
	sY := 0.0

	for i, xp := range x {
		dx := xp - meanX
		dy := y[i] - meanY

		sX += dx * dx
		sY += dy * dy

		sXY += dx * dy
	}

	return sXY / (math.Sqrt(sX) * math.Sqrt(sY))
}

// AbsMax returns max(|x| for x in arr)
func AbsMax(arr []float64) float64 {
	m, _ := ArgAbsMax(arr)
	return m
}

// ArgAbsMax returns max(|x| for x in arr) and the index
func ArgAbsMax(arr []float64) (max float64, idx int) {
	for i, v := range arr {
		if math.Abs(v) > max {
			idx = i
			max = math.Abs(v)
		}
	}
	return max, idx
}

// Avg computes the average of the given array
func Avg(arr []float64) float64 {
	if len(arr) == 0 {
		return 0
	}

	sum := 0.0
	for _, v := range arr {
		sum += v
	}

	return sum / float64(len(arr))
}

// ArgAbove returns the indices for the values from the array that are above the given threshold
func ArgAbove(threshold float64, arr []float64) []int {
	var res []int
	for idx, v := range arr {
		if v > threshold {
			res = append(res, idx)
		}
	}
	return res
}
