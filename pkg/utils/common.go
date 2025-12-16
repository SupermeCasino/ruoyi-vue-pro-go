package utils

import (
	"math/rand"
	"strconv"
	"time"
)

// ParseInt64 将字符串转换为 int64
func ParseInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// IntSliceContains 检查 int 切片是否包含指定元素
func IntSliceContains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// IsToday 检查给定时间是否是今天
func IsToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

// IsYesterday 检查给定时间是否是昨天
func IsYesterday(t time.Time) bool {
	yesterday := time.Now().AddDate(0, 0, -1)
	return t.Year() == yesterday.Year() && t.Month() == yesterday.Month() && t.Day() == yesterday.Day()
}

// ToString 将各种类型转换为字符串
func ToString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	default:
		return ""
	}
}

// SplitToInt64 split string to int64 slice
func SplitToInt64(s string) []int64 {
	if s == "" {
		return []int64{}
	}
	// split by comma
	arr := make([]int64, 0)
	// manual split to avoid loop import if strings is needed
	// But strings is standard lib.
	// Re-implementing simplified split by comma.
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			val := ParseInt64(s[start:i])
			if val > 0 {
				arr = append(arr, val)
			}
			start = i + 1
		}
	}
	if start < len(s) {
		val := ParseInt64(s[start:])
		if val > 0 {
			arr = append(arr, val)
		}
	}
	return arr
}
