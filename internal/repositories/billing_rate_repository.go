package repositories

import (
	interfaces "rentjoy/internal/interfaces/repositories"
	"rentjoy/internal/models"
	"time"

	"gorm.io/gorm"
)

type BillingRateRepository struct {
	*GenericRepository[models.BillingRate]
}

func NewBillingRateRepository(db *gorm.DB) interfaces.BillingRateRepository {
	return &BillingRateRepository{
		GenericRepository: NewGenericRepository[models.BillingRate](db),
	}
}

func (r *BillingRateRepository) FindAvailableTimes(venueID int, dayOfWeek time.Weekday) ([]models.BillingRate, error) {
	var rates []models.BillingRate

	err := r.DB.Where("VenueId = ? AND DayOfWeek = ?", venueID, dayOfWeek).
		Order("StartTime").
		Find(&rates).Error

	return rates, err
}
