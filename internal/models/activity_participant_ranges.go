package models

type ActivityParticipantRange struct {
	ID            uint   `gorm:"primaryKey;column:Id"`
	PeopleOfRange string `gorm:"column:PeopleOfRange"`
}

func (ActivityParticipantRange) TableName() string {
	return "ActivityParticipantRanges"
}
