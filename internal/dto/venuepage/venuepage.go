package venuepage

type VenuePage struct {
	VenueID                 uint          `json:"venueId"`
	ImgUrls                 []string      `json:"imgUrls"`
	Name                    string        `json:"name"`
	City                    string        `json:"city"`
	District                string        `json:"district"`
	Address                 string        `json:"address"`
	NumberOfPeople          string        `json:"numberOfPeople"`
	SpaceSize               string        `json:"spaceSize"`
	VenueDevices            []VenueDevice `json:"venueDevices"`
	VenueNotIncludedDevices []string      `json:"venueNotIncludedDevices"`
	CommentAverage          float32       `json:"commentAverage"`
	VenueComment            []Comment     `json:"venueComment"`
	VenueRules              []string      `json:"venueRules"`
	TrafficInfo             TrafficInfo   `json:"trafficInfo"`
	Lng                     string        `json:"lng"`
	Lat                     string        `json:"lat"`
	OwnerInfo               OwnerInfo     `json:"ownerInfo"`
	UnsubscribeRule         string        `json:"unsubscribeRule"`
	HrPriceRange            string        `json:"hrPriceRange"`
	TimeSlotPriceRange      string        `json:"timeSlotPriceRange"`
	Recommended             []Recommended `json:"recommended"`
	ReservedDate            []string      `json:"reserveDate"`
	OpenDayOfWeek           []int         `json:"openDayOfWeek"`
	MinRentHours            float32       `json:"minRentHours"`
}

type AvailableTime struct {
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
	Price         string `json:"price"`
	BillingRateID string `json:"billingRateId"`
	RateTypeID    string `json:"rateTypeId"`
}

type VenueDevice struct {
	DeviceName     string `json:"deviceName"`
	DeviceQuantity int    `json:"deviceQuantity"`
	DeviceRemark   string `json:"deviceRemark"`
}

type Comment struct {
	UserName     string `json:"userName"`
	CommentYear  string `json:"commentYear"`
	CommentMonth string `json:"commentMonth"`
	CommentDay   string `json:"commentDay"`
	CommentTxt   string `json:"commentTxt"`
}

type TrafficInfo struct {
	MRTInfo  string `json:"mrtInfo"`
	BusInfo  string `json:"busInfo"`
	ParkInfo string `json:"ParkInfo"`
}

type OwnerInfo struct {
	ImgUrl    string `json:"imgUrl"`
	Name      string `json:"name"`
	JoinYear  string `json:"joinYear"`
	JoinMonth string `json:"joinMonth"`
	JoinDay   string `json:"joinDay"`
}
