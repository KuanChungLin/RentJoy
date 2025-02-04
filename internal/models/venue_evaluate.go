package models

import "time"

type VenueEvaluate struct {
	ID              uint      `gorm:"primaryKey;column:Id"`
	CreatedAt       time.Time `gorm:"column:CreateAt"`
	EvaluateRate    uint      `gorm:"column:EvaluateRate"`
	EvaluateComment string    `gorm:"column:EvaluateComment"`

	Order *Order `gorm:"foreignKey:ID"`
}

func (VenueEvaluate) TableName() string {
	return "VenueEvaluate"
}
