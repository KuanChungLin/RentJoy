package venuepage

import "rentjoy/internal/dto/homepage"

// 預定頁
type ReservedPage struct {
	VenueID              string              `json:"venueId"`
	VenueImgUrl          string              `json:"venueImgUrl"`
	Name                 string              `json:"name"`
	Address              string              `json:"address"`
	Date                 string              `json:"date"`
	TimeDetails          []TimeDetail        `json:"timeDetails"`
	Amount               string              `json:"amount"`
	ReservedActivities   []homepage.Activity `json:"reservedActivities"`
	ReservedDetailCookie string              `json:"reservedDetailCookie"`
}

type TimeDetail struct {
	TimeRange string `json:"timeRange"`
	Price     string `json:"price"`
}

// 取回在 JS 存入的 Cookie 資料
// StartTime  為小時制的場地預訂開始時間
// EndTime    為小時制的場地預訂結束時間
// ReservedDay為預訂的日期 (台北標準時間)
type ReservedDetail struct {
	TimeSlotIds []string `json:"timeSlotIds"`
	StartTime   string   `json:"startTime"`
	EndTime     string   `json:"endTime"`
	ReservedDay string   `json:"reservedDay"`
	VenueID     int      `json:"venueId"`
}
