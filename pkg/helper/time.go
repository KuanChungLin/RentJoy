package helper

import (
	"fmt"
	"time"
)

// ParseTime 解析時間字串，返回 time.Time 和錯誤
func ParseTime(timeStr string) (time.Time, error) {
	// 嘗試解析不同的時間格式
	layouts := []string{
		"2006/01/02 15:04:05", // ECPay 常用格式
		"2006-01-02 15:04:05",
		time.RFC3339,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("無法解析時間格式: %s", timeStr)
}

// MustParseTime 解析時間字串，如果失敗返回當前時間
func MustParseTime(timeStr string) time.Time {
	t, err := ParseTime(timeStr)
	if err != nil {
		return time.Now()
	}
	return t
}

// 將時間傳給指針類型的變數
func TimePtr(t time.Time) *time.Time {
	return &t
}
