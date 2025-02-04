package searchpage

import "rentjoy/internal/dto/homepage"

type SearchPage struct {
	ActivityList    []homepage.Activity    `json:"activityList"`
	PeopleCountList []homepage.PeopleCount `json:"peopleCountList"`
	MaxPriceList    []int                  `json:"maxPriceList"`
	MinPriceList    []int                  `json:"minPriceList"`
	VenueInfos      []VenueInfo            `json:"vneueInfos"`
	VenueFilter     VenueFilter            `json:"venueFilter"`
}

type VenueInfo struct {
	VenueImgUrl    string              `json:"venueImgUrl"`
	VenueID        uint                `json:"venueId"`
	VenueName      string              `json:"venueName"`
	VenueOwner     string              `json:"venueOwner"`
	VenueCity      string              `json:"venueCity"`
	VenueDistrict  string              `json:"venueDistrict"`
	NumberOfPeople int                 `json:"numberOfPeople"`
	VenuePrice     string              `json:"venuePrice"`
	ActivityTags   []homepage.Activity `json:"activityTags"`
}

type VenuePartialResponse struct {
	VenueInfos []VenueInfo
	EndOfData  bool
}

type VenueFilter struct {
	ActivityID     uint   `json:"activityId"`
	NumberOfPeople string `json:"numberOfPeople"`
	City           string `json:"city"`
	District       string `json:"district"`
	MaxPrice       int    `json:"maxPrice"`
	MinPrice       int    `json:"minPrice"`
	VenueName      string `json:"venueName"`
	DayType        string `json:"dayType"`
	RentTime       string `json:"rentTime"`
	Page           int    `json:"page"`
}
