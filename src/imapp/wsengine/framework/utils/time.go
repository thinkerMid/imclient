package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandDuration(min time.Duration, max time.Duration) time.Duration {
	return time.Duration(RandInt64(min.Nanoseconds(), max.Nanoseconds()))
}

//func LogRandMillSecond(min time.Duration, max time.Duration) float64 {
//	v := LogRandSecond(min, max) * 1000
//	result, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", v), 64)
//	return result
//}
//
//func LogRandSecond(min time.Duration, max time.Duration) float64 {
//	r := RandDuration(min, max)
//
//	return r.Seconds()
//}

// LogRandSecond 随机范围内的秒
func LogRandSecond(min time.Duration, max time.Duration) float64 {
	minS := min.Seconds()
	maxS := max.Seconds()

	return minS + rand.Float64()*(maxS-minS)
}

// LogRandMillSecond 随机范围毫秒
func LogRandMillSecond(min time.Duration, max time.Duration) float64 {
	minS := min.Seconds()
	maxS := max.Seconds()

	// 转换毫秒单位
	return (minS + rand.Float64()*(maxS-minS)) * 1000
}
