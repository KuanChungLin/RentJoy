package venuepage

import "rentjoy/internal/dto/homepage"

// 預定頁
type ReservedPage struct {
	VenueID            string              `json:"venueId"`
	VenueImgUrl        string              `json:"venueImgUrl"`
	Name               string              `json:"name"`
	Address            string              `json:"address"`
	Date               string              `json:"date"`
	TimeDetails        []TimeDetail        `json:"timeDetails"`
	Amount             string              `json:"amount"`
	ReservedActivities []homepage.Activity `json:"reservedActivities"`
}

type TimeDetail struct {
	TimeRange string `json:"timeRange"`
	Price     string `json:"price"`
}

type ReservedDetail struct {
	TimeSlotIds []string `json:"timeSlotIds"`
	StartTime   string   `json:"startTime"`
	EndTime     string   `json:"endTime"`
	ReservedDay string   `json:"reservedDay"`
	VenueID     int      `json:"venueId"`
}
