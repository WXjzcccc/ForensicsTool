package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// timestampToDatetime 将时间转换为指定时区的格式化字符串
func timestampToDatetime(convertedTime time.Time, originTimezone, targetTimezone string) (string, error) {
	// 加载原始时区
	originLoc, err := time.LoadLocation(originTimezone)
	if err != nil {
		return "", err
	}

	// 加载目标时区
	targetLoc, err := time.LoadLocation(targetTimezone)
	if err != nil {
		return "", err
	}

	// 将时间转换为原始时区
	convertedTime = convertedTime.In(originLoc)
	// 转换为目标时区
	convertedTimeInTarget := convertedTime.In(targetLoc)
	// 格式化时间
	dateStr := convertedTimeInTarget.Format("2006-01-02 15:04:05")
	return dateStr, nil
}

// IosTimestampToDatetime 转换iOS时间戳(从2001-01-01开始的秒数)
func IosTimestampToDatetime(timestamp float64, originTimezone, targetTimezone string) string {
	if originTimezone == "" {
		originTimezone = "UTC"
	}
	if targetTimezone == "" {
		targetTimezone = "Asia/Shanghai"
	}
	baseTime := time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	duration := time.Duration(timestamp * float64(time.Second))
	convertedTime := baseTime.Add(duration)

	result, err := timestampToDatetime(convertedTime, originTimezone, targetTimezone)
	if err != nil {
		return ""
	}
	return result
}

// DefaultTimestampToDatetime 根据时间戳长度自动判断单位并转换
func DefaultTimestampToDatetime(timestamp int64, originTimezone, targetTimezone string) string {
	if originTimezone == "" {
		originTimezone = "UTC"
	}
	if targetTimezone == "" {
		targetTimezone = "Asia/Shanghai"
	}
	tStr := strconv.FormatInt(timestamp, 10)
	tLength := len(tStr)

	var dt time.Time

	switch {
	case tLength == 13: // 毫秒
		dt = time.Unix(timestamp/1000, (timestamp%1000)*int64(time.Millisecond))
	case tLength == 16: // 微秒
		dt = time.Unix(timestamp/1000000, (timestamp%1000000)*int64(time.Microsecond))
	case tLength == 19: // 纳秒
		dt = time.Unix(timestamp/1000000000, timestamp%1000000000)
	default: // 秒
		dt = time.Unix(timestamp, 0)
	}

	result, err := timestampToDatetime(dt, originTimezone, targetTimezone)
	if err != nil {
		return ""
	}
	return result
}

// ChromeTimestampToDatetime 转换Chrome时间戳(从1601-01-01开始的微秒数)
func ChromeTimestampToDatetime(timestamp int64, originTimezone, targetTimezone string) string {
	if originTimezone == "" {
		originTimezone = "UTC"
	}
	if targetTimezone == "" {
		targetTimezone = "Asia/Shanghai"
	}
	baseTime := time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC)
	duration := time.Duration(timestamp) * time.Microsecond
	convertedTime := baseTime.Add(duration)

	result, err := timestampToDatetime(convertedTime, originTimezone, targetTimezone)
	if err != nil {
		return ""
	}
	return result
}

// WindowsFileTimeToDatetime 转换Windows文件时间(从1601-01-01开始的100纳秒间隔)
func WindowsFileTimeToDatetime(timestamp uint64, originTimezone, targetTimezone string) string {
	if originTimezone == "" {
		originTimezone = "UTC"
	}
	if targetTimezone == "" {
		targetTimezone = "Asia/Shanghai"
	}
	baseTime := time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC)
	duration := time.Duration(timestamp/10) * time.Microsecond
	convertedTime := baseTime
	if duration < 0 { //偏移过大，导致溢出变成负值，处理溢出
		maxDuration := 1<<63 - 1
		minDuration := -1 << 63
		convertedTime = convertedTime.Add(time.Duration(maxDuration))
		convertedTime = convertedTime.Add(duration - time.Duration(minDuration))
	} else {
		convertedTime = convertedTime.Add(duration)
	}
	result, err := timestampToDatetime(convertedTime, originTimezone, targetTimezone)
	if err != nil {
		return ""
	}
	return result
}

// NineTimestampToDatetime 转换特殊的9位时间戳
func NineTimestampToDatetime(timestamp int64, originTimezone, targetTimezone string) string {
	if originTimezone == "" {
		originTimezone = "UTC"
	}
	if targetTimezone == "" {
		targetTimezone = "Asia/Shanghai"
	}
	nine := timestamp - 621355968000000000
	return DefaultTimestampToDatetime(nine/10000, originTimezone, targetTimezone)
}

// AppleTimestampToDatetime 转换Apple时间戳
func AppleTimestampToDatetime(timestamp float64, originTimezone, targetTimezone string) string {
	if originTimezone == "" {
		originTimezone = "UTC"
	}
	if targetTimezone == "" {
		targetTimezone = "Asia/Shanghai"
	}
	timestampStr := fmt.Sprintf("%.20f", timestamp)
	parts := strings.Split(timestampStr, ".")
	if len(parts) < 1 {
		return ""
	}

	secondsStr := parts[0]
	if len(secondsStr) < 9 {
		return ""
	}

	seconds, err := strconv.ParseInt(secondsStr[:9], 10, 64)
	if err != nil {
		return ""
	}

	seconds += 978307200 // 从2001-01-01到1970-01-01的秒数
	dt := time.Unix(seconds, 0)

	result, err := timestampToDatetime(dt, originTimezone, targetTimezone)
	if err != nil {
		return ""
	}
	return result
}
