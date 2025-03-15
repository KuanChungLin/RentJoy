package repositories

import (
	"errors"
	"log"
	"rentjoy/internal/dto/order"
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

func (r *OrderRepository) FindByUserAndStatus(userId uint, status order.OrderStatus, pageIndex int, pageSize int) ([]models.Order, error) {
	var orders []models.Order
	var err error

	// 計算要跳過的記錄數
	offset := (pageIndex - 1) * pageSize

	// 執行查詢，包含分頁和降序排序
	if status == 3 {
		err = r.DB.Where("MemberId = ? AND OrderStatus = ?", userId, status).
			Order("UnsubscribeTime DESC").
			Offset(offset).
			Limit(pageSize).
			Find(&orders).Error
	} else {
		err = r.DB.Where("MemberId = ? AND OrderStatus = ?", userId, status).
			Order("CreateAt DESC").
			Offset(offset).
			Limit(pageSize).
			Find(&orders).Error
	}

	return orders, err
}

// 透過 userId 及 orderStatus 取得使用者的預訂單數量
func (r *OrderRepository) CountByUserAndStatus(userId uint, status order.OrderStatus) (int, error) {
	var count int64

	err := r.DB.Model(&models.Order{}).
		Where("MemberId = ? AND OrderStatus = ?", userId, status).
		Count(&count).Error

	return int(count), err
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
		Where("Id = ?", id).
		Update("Status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("order not found")
	}

	return nil
}

// 透過 VenueId 取得所有訂單
func (r *OrderRepository) FindOrdersByVenueId(tx *gorm.DB, venueId uint) ([]models.Order, error) {
	var orders []models.Order

	err := tx.Where("VenueId = ?", venueId).
		Preload("VenueEvaluate").
		Find(&orders).Error
	if err != nil {
		log.Printf("Find Orders By VenueId Error:%s", err)
		return orders, err
	}

	return orders, err
}

// 透過 UserId 取得所有預訂單
func (r *OrderRepository) FindManageOrderByUserId(userId uint) ([]models.Order, error) {
	var orders []models.Order

	err := r.DB.Where("VenueOwnerId = ?", userId).
		Order("CreateAt Desc").
		Preload("Venue").
		Preload("Details").
		Find(&orders).Error
	if err != nil {
		log.Printf("Find ManageOrders By UserId Error:%s", err)
		return orders, err
	}

	return orders, err
}
