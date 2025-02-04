package models

import "time"

type Management struct {
	ID                    uint      `gorm:"primaryKey;column:Id"`
	MemberID              uint      `gorm:"column:MemberId"`
	ManagementName        string    `gorm:"column:ManagementName"`
	ManagementDescription string    `gorm:"column:ManagementDescription"`
	AvatarImgLinkPath     string    `gorm:"column:AvatarImgLinkPath"`
	Contact               string    `gorm:"column:Contact"`
	PublicNumber          string    `gorm:"column:PublicNumber"`
	PrivateNumber         string    `gorm:"column:PrivateNumber"`
	CreatedAt             time.Time `gorm:"column:CreateAt"`
	UpdatedAt             time.Time `gorm:"column:EditAt"`
	IsDeleted             bool      `gorm:"column:IsDelete"`

	Member Member `gorm:"foreignKey:MemberID"`

	ManagedVenues []VenueInformation `gorm:"foreignKey:ManagementID"`
}

func (Management) TableName() string {
	return "Management"
}
