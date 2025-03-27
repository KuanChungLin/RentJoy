package interfaces

import "rentjoy/internal/dto/create"

type CreateService interface {
	GetSpaceTypes() ([]create.SpaceType, error)
	GetActivities() ([]create.ActivityInfo, error)
	GetEquipmentTypes() ([]create.EquipmentType, error)
	GetManagers(userId uint) ([]create.ManagerInfo, error)
	CreateVenue(userId uint, form *create.CreateForm) error
}
