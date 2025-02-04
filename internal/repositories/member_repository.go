package repositories

import (
	interfaces "rentjoy/internal/interfaces/repositories"
	"rentjoy/internal/models"

	"gorm.io/gorm"
)

type MemberRepository struct {
	*GenericRepository[models.Member]
}

func NewMemberRepository(db *gorm.DB) interfaces.MemberRepository {
	return &MemberRepository{
		GenericRepository: NewGenericRepository[models.Member](db),
	}
}

func (r *MemberRepository) FindByID(id uint) (*models.Member, error) {
	var member models.Member
	err := r.DB.
		Preload("FacebookLogins"). // 預加載 Facebook 登入資料
		Preload("GoogleLogins").   // 預加載 Google 登入資料
		First(&member, id).Error

	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *MemberRepository) FindByAccount(a string) (*models.Member, error) {
	var m models.Member
	err := r.GenericRepository.DB.Where("Account = ?", a).
		First(&m).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &m, err
}

func (r *MemberRepository) IsEmailExists(email string, exceptUserID uint) (bool, error) {
	var count int64
	err := r.DB.Model(&models.Member{}).
		Where("email = ? AND id != ?", email, exceptUserID).
		Count(&count).Error

	return count > 0, err
}

func (r *MemberRepository) FindByFacebookID(ID string) (*models.Member, error) {
	var m models.Member
	err := r.GenericRepository.DB.Where("IsDelete = ?", 0).
		Joins("JOIN FacebookThirPartyLogin fb ON fb.MemberId = Members.Id").
		Where("fb.FacebookThirPartyId = ?", ID).
		Find(&m).Error

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *MemberRepository) FindByGoogleID(ID string) (*models.Member, error) {
	var m models.Member
	err := r.GenericRepository.DB.Where("IsDelete = ?", 0).
		Joins("JOIN GoogleThirPartyLogin g ON g.MemberId = Members.Id").
		Where("g.GoogleThirPartyId = ?", ID).
		Find(&m).Error

	if err != nil {
		return nil, err
	}

	return &m, nil
}
