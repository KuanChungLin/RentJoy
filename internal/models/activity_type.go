package models

type ActivityType struct {
	ID                  uint   `gorm:"primaryKey;column:Id"`
	ActivityName        string `gorm:"column:ActivityName"`
	ActivityIcon        string `gorm:"column:ActivityIcon"`
	ActivityDescription string `gorm:"column:ActivityDescription"`

	Orders     []Order         `gorm:"foreignKey:ActivityTypeID"`
	Activities []VenueActivity `gorm:"foreignKey:ActivityID"`
}

func (ActivityType) TableName() string {
	return "ActivityType"
}
