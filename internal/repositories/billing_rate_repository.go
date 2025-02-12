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

func (r *BillingRateRepository) FindByReserved(venueID uint, rateTypeID uint, weekDay int) (*models.BillingRate, error) {
	var rate models.BillingRate

	err := r.DB.Where("VenueId = ? AND RateTypeId = ? AND DayOfWeek = ?", venueID, rateTypeID, weekDay).
		First(&rate).Error

	return &rate, err
}
