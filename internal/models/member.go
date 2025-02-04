package models

import "time"

type Member struct {
	ID        uint      `gorm:"primaryKey;column:Id"`
	FirstName string    `gorm:"column:FirstName"`
	LastName  string    `gorm:"column:LastName"`
	Account   string    `gorm:"column:Account"`
	Password  string    `gorm:"column:Password"`
	Email     string    `gorm:"column:Email"`
	Phone     string    `gorm:"column:Phone"`
	IsDeleted bool      `gorm:"column:IsDelete"`
	CreatedAt time.Time `gorm:"column:CreateAt"`
	UpdatedAt time.Time `gorm:"column:EditAt"`

	FacebookLogins    []FacebookThirdPartyLogin `gorm:"foreignKey:MemberID"`
	GoogleLogins      []GoogleThirdPartyLogin   `gorm:"foreignKey:MemberID"`
	Managements       []Management              `gorm:"foreignKey:MemberID"`
	OrderAsMember     []Order                   `gorm:"foreignKey:MemberID"`
	OrderAsVenueOwner []Order                   `gorm:"foreignKey:OwnerID"`
	OwnedVenues       []VenueInformation        `gorm:"foreignKey:OwnerID"`
}

func (Member) TableName() string {
	return "Members"
}
