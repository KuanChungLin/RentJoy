package models

type DeviceType struct {
	ID             uint   `gorm:"primaryKey;column:Id"`
	DeviceTypeName string `gorm:"column:DeviceTypeName"`

	DeviceItems []DeviceItem `gorm:"foreignKey:DeviceTypeID"`
}

func (DeviceType) TableName() string {
	return "DeviceType"
}
