package utils

import "math"

func DollarsToCents(d float64) int64 {
	return int64(math.Round(d * 100))
}

func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func RateToBasisPoint(r float64) int64 {
	return int64(math.Round(r * 10000))
}

func DivRoundHalfUp(n, div int64) int64 {
	return (n + div/2) / div
}
