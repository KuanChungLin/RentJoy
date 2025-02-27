package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID              uint            `gorm:"primaryKey;column:Id"`
	Status          int             `gorm:"column:OrderStatus"`
	VenueID         uint            `gorm:"column:VenueId"`
	ActivityTypeID  uint            `gorm:"column:OrderForActivityType"`
	FirstName       string          `gorm:"column:SubscriberFirstName"`
	LastName        string          `gorm:"column:SubscriberLastName"`
	Phone           string          `gorm:"column:SubscriberPhone"`
	Email           string          `gorm:"column:SubscriberEmail"`
	Amount          decimal.Decimal `gorm:"column:Amount;type:datetime"`
	CreatedAt       time.Time       `gorm:"column:CreateAt"`
	UnsubscribeTime *time.Time      `gorm:"column:UnsubscribeTime"` // 使用指針方可傳入 nil
	UserCount       int             `gorm:"column:UserCount"`
	Message         string          `gorm:"column:UserMessage"`
	MemberID        uint            `gorm:"column:MemberId"`
	OwnerID         uint            `gorm:"column:VenueOwnerId"`

	EcpayOrder    *EcpayOrder      `gorm:"foreignKey:ID"`
	Member        Member           `gorm:"foreignKey:MemberID"`
	ActivityType  ActivityType     `gorm:"foreignKey:ActivityTypeID"`
	Venue         VenueInformation `gorm:"foreignKey:VenueID"`
	VenueOwner    Member           `gorm:"foreignKey:OwnerID"`
	VenueEvaluate *VenueEvaluate   `gorm:"foreignKey:ID"`

	Details []OrderDetail `gorm:"foreignKey:OrderID"`
}

func (Order) TableName() string {
	return "Order"
}
