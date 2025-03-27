package repositories

import (
	interfaces "rentjoy/internal/interfaces/repositories"
	"rentjoy/internal/models"

	"gorm.io/gorm"
)

type SpaceTypeRepository struct {
	*GenericRepository[models.SpaceType]
}

func NewSpaceTypeRepository(db *gorm.DB) interfaces.SpaceTypeRepository {
	return &SpaceTypeRepository{
		GenericRepository: NewGenericRepository[models.SpaceType](db),
	}
}

// 取得所有場地、空間類型
func (r *SpaceTypeRepository) FindSpaceAndFacility() ([]models.SpaceType, error) {
	var spaceTypes []models.SpaceType

	err := r.DB.Preload("VenueTypes").Find(&spaceTypes).Error
	if err != nil {
		return nil, err
	}

	return spaceTypes, nil
}
