package repositories

import (
	interfaces "rentjoy/internal/interfaces/repositories"
	"rentjoy/internal/models"

	"gorm.io/gorm"
)

type VenueImgRepository struct {
	*GenericRepository[models.VenueImg]
}

func NewVenueImgRepository(db *gorm.DB) interfaces.VenueImgRepository {
	return &VenueImgRepository{
		GenericRepository: NewGenericRepository[models.VenueImg](db),
	}
}

func (r *VenueImgRepository) FindFirstBySort(venueID uint, sort int) (*models.VenueImg, error) {
	var img models.VenueImg

	err := r.DB.Where("VenueId = ? AND Sort = ?", venueID, sort).
		First(&img).Error

	if err != nil {
		return nil, err
	}

	return &img, nil
}
