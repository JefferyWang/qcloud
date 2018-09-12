package util

import (
	"strconv"
	"time"
)

// GetCurTimeStampStr 获取当前时间戳字符串
func GetCurTimeStampStr() string {
	timestamp := time.Now().Unix()
	return strconv.FormatInt(timestamp, 10)
}
