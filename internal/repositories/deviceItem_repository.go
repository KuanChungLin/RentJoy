package repositories

import (
	interfaces "rentjoy/internal/interfaces/repositories"
	"rentjoy/internal/models"

	"gorm.io/gorm"
)

type DeviceItemRepository struct {
	*GenericRepository[models.DeviceItem]
}

func NewDeviceItemRepository(db *gorm.DB) interfaces.DeviceItemRepository {
	return &DeviceItemRepository{
		GenericRepository: NewGenericRepository[models.DeviceItem](db),
	}
}

func (r *DeviceItemRepository) GetAllDeviceItemNames() ([]string, error) {
	var deviceItems []string

	err := r.DB.Model(&models.DeviceItem{}).
		Pluck("DeviceName", &deviceItems).Error

	if err != nil {
		return nil, err
	}

	return deviceItems, err
}
