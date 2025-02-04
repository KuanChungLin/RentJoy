package models

type FacebookThirdPartyLogin struct {
	ID                   uint   `gorm:"primaryKey;column:Id"`
	MemberID             uint   `gorm:"column:MemberId"`
	FacebookThirdPartyID string `gorm:"column:FacebookThirPartyId"`

	Member Member `gorm:"foreignKey:MemberID"`
}

func (FacebookThirdPartyLogin) TableName() string {
	return "FacebookThirPartyLogin"
}
