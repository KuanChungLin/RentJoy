package repositories

import (
	interfaces "rentjoy/internal/interfaces/repositories"
	"rentjoy/internal/models"

	"gorm.io/gorm"
)

type GoogleRepository struct {
	*GenericRepository[models.GoogleThirdPartyLogin]
}

func NewGoogleRepository(db *gorm.DB) interfaces.GoogleRepository {
	return &GoogleRepository{
		GenericRepository: NewGenericRepository[models.GoogleThirdPartyLogin](db),
	}
}