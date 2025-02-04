package repositories

import (
	interfaces "rentjoy/internal/interfaces/repositories"
	"rentjoy/internal/models"

	"gorm.io/gorm"
)

type ParticipantRangeRepository struct {
	*GenericRepository[models.ActivityParticipantRange]
}

func NewParticipantRangeRepository(db *gorm.DB) interfaces.ParticipantRangeRepository {
	return &ParticipantRangeRepository{
		GenericRepository: NewGenericRepository[models.ActivityParticipantRange](db),
	}
}
