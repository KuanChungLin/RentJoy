package repositories

import (
	"errors"
	interfaces "rentjoy/internal/interfaces/repositories"
	"rentjoy/internal/models"

	"gorm.io/gorm"
)

type VenueEvaluateRepository struct {
	*GenericRepository[models.VenueEvaluate]
}

func NewVenueEvaluateRepository(db *gorm.DB) interfaces.VenueEvaluateRepository {
	return &VenueEvaluateRepository{
		GenericRepository: NewGenericRepository[models.VenueEvaluate](db),
	}
}

func (r *VenueEvaluateRepository) FindByOrderId(orderId uint) (*models.VenueEvaluate, error) {
	var evaluate *models.VenueEvaluate

	err := r.DB.Where("Id = ?", orderId).
		First(evaluate).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 找不到记录时返回 nil, nil
		}
		return nil, err
	}
	return evaluate, nil
}
