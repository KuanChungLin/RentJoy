package repositories

import (
	"errors"
	"log"
	interfaces "rentjoy/internal/interfaces/repositories"
	"rentjoy/internal/models"

	"gorm.io/gorm"
)

type VenueEvaluateRepository struct {
	*GenericRepository[models.VenueEvaluate]
}

func NewVenueEvaluateRepository(db *gorm.DB) interfaces.VenueEvaluateRepository {
	return &VenueEvaluateRepository{
		GenericRepository: NewGenericRepository[models.VenueEvaluate](db),
	}
}

// 透過 OrderId 取得訂單評價
func (r *VenueEvaluateRepository) FindByOrderId(orderId uint) (*models.VenueEvaluate, error) {
	var evaluate models.VenueEvaluate

	err := r.DB.Where("Id = ?", orderId).
		Find(&evaluate).Error
	if err != nil {
		log.Println(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 找不到记录时返回 nil, nil
		}
		return nil, err
	}

	return &evaluate, nil
}

// 交易流程建立場地評價資料
func (r *VenueEvaluateRepository) CreateByTx(tx *gorm.DB, evaluate models.VenueEvaluate) error {
	err := tx.Create(&evaluate).Error
	if err != nil {
		log.Printf("Create VenueEvaluate Error:%s", err)
		return err
	}
	return nil
}
