package models

type VenueDevice struct {
	ID                uint   `gorm:"primaryKey;column:Id"`
	DeviceItemID      uint   `gorm:"column:DeviceItemId"`
	VenueID           uint   `gorm:"column:VenueId"`
	Count             int    `gorm:"column:Count"`
	DeviceDescription string `gorm:"column:DeviceDescription"`

	DeviceItem DeviceItem       `gorm:"foreignKey:DeviceItemID"`
	Venue      VenueInformation `gorm:"foreignKey:VenueID"`
}

func (VenueDevice) TableName() string {
	return "VenueDevices"
}
