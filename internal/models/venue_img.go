package models

type VenueImg struct {
	ID           uint   `gorm:"primaryKey;column:Id"`
	VenueID      uint   `gorm:"column:VenueId"`
	VenueImgPath string `gorm:"column:VenueImgPath"`
	Sort         uint   `gorm:"column:Sort"`
	PublicID     string `gorm:"column:PublicID"`

	Venue VenueInformation `gorm:"foreignKey:VenueID"`
}

func (VenueImg) TableName() string {
	return "VenueImg"
}
