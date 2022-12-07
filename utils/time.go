package utils

import "time"

// 获取时间戳
func GetCurTime() int64 {
	timeUnix := time.Now().Unix()
	return timeUnix
}
