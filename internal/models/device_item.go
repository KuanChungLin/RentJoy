package models

type DeviceItem struct {
	ID           uint   `gorm:"primaryKey;column:Id"`
	DeviceTypeID uint   `gorm:"column:DeviceTypeId"`
	DeviceName   string `gorm:"column:DeviceName"`

	DeviceType   DeviceType    `gorm:"foreignKey:DeviceTypeID"`
	VenueDevices []VenueDevice `gorm:"foreignKey:DeviceItemID"`
}

func (DeviceItem) TableName() string {
	return "DeviceItems"
}
