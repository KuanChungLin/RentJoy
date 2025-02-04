package repositories

import (
	interfaces "rentjoy/internal/interfaces/repositories"
	"rentjoy/internal/models"

	"gorm.io/gorm"
)

type ActivityTypeRepository struct {
	*GenericRepository[models.ActivityType]
}

func NewActivityTypeRepository(db *gorm.DB) interfaces.ActivityTypeRepository {
	return &ActivityTypeRepository{
		GenericRepository: NewGenericRepository[models.ActivityType](db),
	}
}
