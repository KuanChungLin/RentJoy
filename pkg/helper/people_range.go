package helper

import (
	"math"
)

// 定義範圍結構
type PeopleRange struct {
	Min int
	Max int
}

// 定義人數範圍對應表
var NumberOfPeopleRanges = map[string]PeopleRange{
	"1 - 10":    {Min: 1, Max: 10},
	"11 - 20":   {Min: 11, Max: 20},
	"21 - 40":   {Min: 21, Max: 40},
	"41 - 60":   {Min: 41, Max: 60},
	"61 - 80":   {Min: 61, Max: 80},
	"81 - 100":  {Min: 81, Max: 100},
	"101 - 200": {Min: 101, Max: 200},
	"201 - 300": {Min: 201, Max: 300},
	"301 - 400": {Min: 301, Max: 400},
	"401 - 500": {Min: 401, Max: 500},
	"500+":      {Min: 501, Max: math.MaxInt32},
}

// 取得人數過濾條件
func GetNumberOfPeopleFilter(rangeStr string) (max, min int) {
	if rangeData, ok := NumberOfPeopleRanges[rangeStr]; ok {
		return rangeData.Max, rangeData.Min
	}
	return 0, 0
}
