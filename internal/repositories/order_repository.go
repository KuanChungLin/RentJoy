package repositories

import (
	interfaces "rentjoy/internal/interfaces/repositories"
	"rentjoy/internal/models"
	"time"

	"gorm.io/gorm"
)

type OrderRepository struct {
	*GenericRepository[models.Order]
}

func NewOrderRepository(db *gorm.DB) interfaces.OrderRepository {
	return &OrderRepository{
		GenericRepository: NewGenericRepository[models.Order](db),
	}
}

func (r *OrderRepository) FindConflictingOrders(venueID int, date time.Time) ([]models.Order, error) {
	var orders []models.Order

	err := r.DB.Where("VenueId = ? AND OrderStatus >= ? AND OrderStatus <= ?", venueID, 0, 2).
		Preload("Details", "DATE(StartTime) = ?", date.Format("2006-01-02")).
		Find(&orders).Error

	return orders, err
}
