package repositories

import (
	"errors"
	interfaces "rentjoy/internal/interfaces/repositories"
	"rentjoy/internal/models"

	"gorm.io/gorm"
)

type EcpayRepository struct {
	*GenericRepository[models.EcpayOrder]
}

func NewEcpayRepository(db *gorm.DB) interfaces.EcpayRepository {
	return &EcpayRepository{
		GenericRepository: NewGenericRepository[models.EcpayOrder](db),
	}
}

func (r *EcpayRepository) CreateEcpayOrder(tx *gorm.DB, ecpayOrder models.EcpayOrder) error {
	if err := tx.Create(&ecpayOrder).Error; err != nil {
		return err
	}

	return nil
}

func (r *EcpayRepository) FindByMerchantTradeNo(tradeNo string) (*models.EcpayOrder, error) {
	var order models.EcpayOrder
	err := r.DB.Where("MerchantTradeNo = ?", tradeNo).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 找不到資料返回 nil
		}
		return nil, err // 其他錯誤返回錯誤
	}
	return &order, nil
}

func (r *EcpayRepository) UpdateByTx(tx *gorm.DB, ecpayOrder models.EcpayOrder) error {
	return tx.Save(&ecpayOrder).Error
}
