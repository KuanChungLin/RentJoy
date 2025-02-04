package models

type RateType struct {
	ID              uint   `gorm:"primaryKey;column:Id"`
	RateDescription string `gorm:"column:RateDescription"`

	BillingRates []BillingRate `gorm:"foreignKey:RateTypeID"`
}

func (RateType) TableName() string {
	return "RateType"
}
