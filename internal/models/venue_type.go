package models

type VenueType struct {
	ID       uint   `gorm:"primaryKey;column:Id"`
	SpaceID  uint   `gorm:"column:SpaceId"`
	TypeName string `gorm:"column:VenueTypeName"`

	Venues []VenueInformation `gorm:"foreignKey:VenueTypeID"`

	SpaceType SpaceType `gorm:"foreignKey:SpaceID"`
}

func (VenueType) TableName() string {
	return "VenueType"
}
