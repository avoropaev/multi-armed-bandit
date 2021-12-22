package bandit

import (
	"math"
)

func sum(values ...int) int {
	var total int

	for _, v := range values {
		total += v
	}

	return total
}

func max(values ...float64) (index int) {
	value := math.Inf(-1)

	for i, v := range values {
		if v > value {
			value = v
			index = i
		}
	}

	return
}
