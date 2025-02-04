package models

import "time"

type VenueInformation struct {
	ID              uint      `gorm:"primaryKey;column:Id"`
	VenueTypeID     uint      `gorm:"column:VenueType"`
	ManagementID    uint      `gorm:"column:ManagementId"`
	Name            string    `gorm:"column:VenueName"`
	Rules           string    `gorm:"column:VenueRules"`
	UnsubscribeRule string    `gorm:"column:UnsubscribeRule"`
	Address         string    `gorm:"column:VenueAddress"`
	Latitude        string    `gorm:"column:Latitude"`
	Status          uint      `gorm:"column:Status"`
	City            string    `gorm:"column:City"`
	District        string    `gorm:"column:District"`
	MRTInfo         string    `gorm:"column:MRTInfo"`
	BusInfo         string    `gorm:"column:BusInfo"`
	ParkInfo        string    `gorm:"column:ParkInfo"`
	NumOfPeople     int       `gorm:"column:NumberOfPeople"`
	SpaceSize       int       `gorm:"column:SpaceSize"`
	CreatedAt       time.Time `gorm:"column:CreateAt"`
	UpdatedAt       time.Time `gorm:"column:EditAt"`
	EvaluateRate    uint      `gorm:"column:AvgEvaluateRate"`
	Longitude       string    `gorm:"column:Longitude"`
	OwnerID         uint      `gorm:"column:VenueOwnerId"`

	BillingRates []BillingRate   `gorm:"foreignKey:VenueID"`
	Orders       []Order         `gorm:"foreignKey:VenueID"`
	Activities   []VenueActivity `gorm:"foreignKey:VenueID"`
	Devices      []VenueDevice   `gorm:"foreignKey:VenueID"`
	Imgs         []VenueImg      `gorm:"foreignKey:VenueID"`

	Management Management `gorm:"foreignKey:ManagementID"`
	Owner      Member     `gorm:"foreignKey:OwnerID"`
	VenueType  VenueType  `gorm:"foreignKey:VenueTypeID"`
}

func (VenueInformation) TableName() string {
	return "VenueInformations"
}
