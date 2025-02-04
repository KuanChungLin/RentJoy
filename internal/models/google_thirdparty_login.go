package models

type GoogleThirdPartyLogin struct {
	ID                 uint   `gorm:"primaryKey;column:Id"`
	MemberID           uint   `gorm:"column:MemberId"`
	GoogleThirdPartyID string `gorm:"column:GoogleThirPartyId"`

	Member Member `gorm:"foreignKey:MemberID"`
}

func (GoogleThirdPartyLogin) TableName() string {
	return "GoogleThirPartyLogin"
}
