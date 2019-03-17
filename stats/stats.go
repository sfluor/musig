package stats

import "math"

// AbsMax returns max(|x| for x in arr)
func AbsMax(arr []float64) (max float64) {
	for _, v := range arr {
		if math.Abs(v) > max {
			max = math.Abs(v)
		}
	}
	return max
}
