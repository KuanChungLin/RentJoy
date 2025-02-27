package repositories

import (
	"errors"
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

// 透過 Service 交易流程創建訂單
func (r *OrderRepository) CreateOrder(tx *gorm.DB, order *models.Order) error {
	return tx.Create(&order).Error
}

// 取得指定場地且狀態未結束的訂單資料
func (r *OrderRepository) FindConflictingOrders(venueID int, date time.Time) ([]models.Order, error) {
	var orders []models.Order

	err := r.DB.Where("VenueId = ? AND OrderStatus >= ? AND OrderStatus <= ?", venueID, 0, 2).
		Preload("Details", "CONVERT(DATE, StartTime) = ?", date.Format("2006-01-02")).
		Find(&orders).Error

	return orders, err
}

// 透過 EcpayID 取得訂單資料
func (r *OrderRepository) FindByEcpayID(tx *gorm.DB, id uint) (*models.Order, error) {
	var order *models.Order

	err := r.DB.Where("Id = ?", id).Find(&order).Error

	return order, err
}

// 更新訂單狀態
func (r *OrderRepository) UpdateStatus(tx *gorm.DB, id uint, status int) error {
	result := tx.Model(&models.Order{}).
		Where("id = ?", id).
		Update("Status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("order not found")
	}

	return nil
}
