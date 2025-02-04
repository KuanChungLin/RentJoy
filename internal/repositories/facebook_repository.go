package repositories

import (
	interfaces "rentjoy/internal/interfaces/repositories"
	"rentjoy/internal/models"

	"gorm.io/gorm"
)

type FacebookRepository struct {
	*GenericRepository[models.FacebookThirdPartyLogin]
}

func NewFacebookRepository(db *gorm.DB) interfaces.FacebookRepository {
	return &FacebookRepository{
		GenericRepository: NewGenericRepository[models.FacebookThirdPartyLogin](db),
	}
}
