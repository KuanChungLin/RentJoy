package helper

import "time"

// 檢查時間是否重疊
func HasTimeOverlap(start1, end1, start2, end2 time.Time) bool {
	// 如果時段1的開始時間等於時段2的結束時間，或時段1的結束時間等於時段2的開始時間，則不算重疊
	return start1.Before(end2) && end1.After(start2) && !start1.Equal(end2) && !end1.Equal(start2)
}

// 組合日期和時間
func CombineDateTime(date time.Time, startTime, endTime time.Time) struct {
	StartTime time.Time
	EndTime   time.Time
} {
	return struct {
		StartTime time.Time
		EndTime   time.Time
	}{
		StartTime: time.Date(
			date.Year(), date.Month(), date.Day(),
			startTime.Hour(), startTime.Minute(), startTime.Second(),
			0, time.UTC,
		),
		EndTime: time.Date(
			date.Year(), date.Month(), date.Day(),
			endTime.Hour(), endTime.Minute(), endTime.Second(),
			0, time.UTC,
		),
	}
}
