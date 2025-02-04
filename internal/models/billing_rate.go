package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type BillingRate struct {
	ID           uint            `gorm:"primaryKey;column:Id"`
	VenueID      uint            `gorm:"column:VenueId"`
	RateTypeID   uint            `gorm:"column:RateTypeId"`
	DayOfWeek    time.Weekday    `gorm:"column:DayOfWeek"`
	StartTime    time.Time       `gorm:"column:StartTime"`
	EndTime      time.Time       `gorm:"column:EndTime"`
	Rate         decimal.Decimal `gorm:"column:Rate"`
	MinRentHours float32         `gorm:"column:MinRentHours"`

	RateType RateType         `gorm:"foreignKey:RateTypeID"`
	Venue    VenueInformation `gorm:"foreignKey:VenueID"`
}

func (BillingRate) TableName() string {
	return "BillingRate"
}
