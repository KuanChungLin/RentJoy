package models

type SpaceType struct {
	ID       uint   `gorm:"primaryKey;column:Id"`
	TypeName string `gorm:"column:SpaceTypeName"`

	VenueTypes []VenueType `gorm:"foreignKey:SpaceID"`
}

func (SpaceType) TableName() string {
	return "SpaceType"
}
