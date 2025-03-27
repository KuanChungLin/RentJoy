package repositories

import (
	interfaces "rentjoy/internal/interfaces/repositories"
	"rentjoy/internal/models"

	"gorm.io/gorm"
)

type ManagementRepository struct {
	*GenericRepository[models.Management]
}

func NewManagementRepository(db *gorm.DB) interfaces.ManagementRepository {
	return &ManagementRepository{
		GenericRepository: NewGenericRepository[models.Management](db),
	}
}

func (r *ManagementRepository) FindByUserId(userId uint) ([]models.Management, error) {
	var managements []models.Management

	err := r.DB.Where("MemberId = ?", userId).Find(&managements).Error
	if err != nil {
		return nil, err
	}

	return managements, nil
}
