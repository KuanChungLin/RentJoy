package create

import "mime/multipart"

type CreateForm struct {
	// 基本資料
	VenueTypeId        int    `json:"facilityId"`
	SelectedActivities []int  `json:"selectedActivities"`
	VenueName          string `json:"venueName"`
	VenueRule          string `json:"venueRule"`
	UnsubscribeRule    string `json:"unsubscribeRule"`

	// 位置資訊
	CitySelect            string `json:"citySelect"`
	DistrictSelect        string `json:"districtSelect"`
	Address               string `json:"areaInfoStreetAddress"`
	TransportationMRT     string `json:"transportInfoMRT"`
	TransportationBus     string `json:"transportInfoBus"`
	TransportationParking string `json:"transportInfoParking"`

	// 場地配置
	SpaceSize      int `json:"spaceSize"`
	NumberOfPeople int `json:"numberOfSpace"`

	// 設備
	Equipments []Equipment `json:"equipmentInfos"`

	// 價格設置
	HourPricing   HourPricingConfig   `json:"hourPricing"`
	PeriodPricing PeriodPricingConfig `json:"periodPricing"`

	// 管理者資訊
	ManagerID int `json:"managerId"`

	// 相片
	VenueImgs []*multipart.FileHeader `json:"venueImgs"`
}

type Equipment struct {
	Id          int    `json:"Id"`
	Quantity    string `json:"Quantity"`
	Description string `json:"Description"`
}

type HourPricingConfig struct {
	LeastRentHours  float32          `json:"leastRentHours"`
	PricingSettings []PricingSetting `json:"pricingSettings"`
}

type PeriodPricingConfig struct {
	PricingSettings []PricingSetting `json:"pricingSettings"`
}

type PricingSetting struct {
	Day       int     `json:"day"`
	StartTime string  `json:"startTime"`
	EndTime   string  `json:"endTime"`
	Price     float64 `json:"price"`
}
