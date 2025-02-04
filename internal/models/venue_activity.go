package models

type VenueActivity struct {
	ID         uint `gorm:"primaryKey;column:Id"`
	ActivityID uint `gorm:"column:ActivityId"`
	VenueID    uint `gorm:"column:VenueId"`

	Activity ActivityType     `gorm:"foreignKey:ActivityID"`
	Venue    VenueInformation `gorm:"foreignKey:VenueID"`
}

func (VenueActivity) TableName() string {
	return "VenueActivities"
}
