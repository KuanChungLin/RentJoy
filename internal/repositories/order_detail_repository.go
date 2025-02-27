package repositories

import (
	"rentjoy/internal/dto/venuepage"
	interfaces "rentjoy/internal/interfaces/repositories"
	"rentjoy/internal/models"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrderDetailRepository struct {
	*GenericRepository[models.OrderDetail]
}

func NewOrderDetailRepository(db *gorm.DB) interfaces.OrderDetailRepository {
	return &OrderDetailRepository{
		GenericRepository: NewGenericRepository[models.OrderDetail](db),
	}
}

// 同步 Order 產生 OrderDetail
func (r *OrderDetailRepository) CreateOrderDetails(tx *gorm.DB, orderID uint, timeDetail *venuepage.ReservedDetail, priceList []decimal.Decimal) error {
	if timeDetail.StartTime == "" {
		// 時段制
		var billingRates []models.BillingRate
		if err := tx.Where("Id IN ?", timeDetail.TimeSlotIds).Find(&billingRates).Error; err != nil {
			return err
		}

		day, err := time.Parse(time.RFC3339, timeDetail.ReservedDay)
		if err != nil {
			return err
		}

		for i, rate := range billingRates {
			orderDetail := models.OrderDetail{
				OrderID: orderID,
				StartTime: time.Date(day.Year(), day.Month(), day.Day(),
					rate.StartTime.Hour(), rate.StartTime.Minute(), rate.StartTime.Second(),
					0, time.UTC),
				EndTime: time.Date(day.Year(), day.Month(), day.Day(),
					rate.EndTime.Hour(), rate.EndTime.Minute(), rate.EndTime.Second(),
					0, time.UTC),
				Price: priceList[i],
			}

			if err := tx.Create(&orderDetail).Error; err != nil {
				return err
			}
		}
		return nil
	}

	// 小時制
	startTime, err := time.Parse(time.RFC3339, timeDetail.StartTime)
	if err != nil {
		return err
	}
	endTime, err := time.Parse(time.RFC3339, timeDetail.EndTime)
	if err != nil {
		return err
	}

	orderDetail := models.OrderDetail{
		OrderID:   orderID,
		StartTime: startTime,
		EndTime:   endTime,
		Price:     priceList[0],
	}

	return tx.Create(&orderDetail).Error
}
