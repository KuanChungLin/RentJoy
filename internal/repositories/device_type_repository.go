package repositories

import (
	interfaces "rentjoy/internal/interfaces/repositories"
	"rentjoy/internal/models"

	"gorm.io/gorm"
)

type DeviceTypeRepository struct {
	*GenericRepository[models.DeviceType]
}

func NewDeviceTypeRepository(db *gorm.DB) interfaces.DeviceTypeRepository {
	return &DeviceTypeRepository{
		GenericRepository: NewGenericRepository[models.DeviceType](db),
	}
}

func (r *DeviceTypeRepository) FindTypeAndItems() ([]models.DeviceType, error) {
	var deviceTypes []models.DeviceType

	err := r.DB.Preload("DeviceItems").Find(&deviceTypes).Error
	if err != nil {
		return nil, err
	}

	return deviceTypes, nil
}
