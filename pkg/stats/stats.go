package stats

import "math"

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
	res := []int{}
	for idx, v := range arr {
		if v > threshold {
			res = append(res, idx)
		}
	}
	return res
}
