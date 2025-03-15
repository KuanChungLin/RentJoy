package interfaces

import (
	"rentjoy/internal/dto/manage"
)

type ManageService interface {
	GetReservedManagement(userId uint) (*manage.ReservedManagement, error)
	ReservedAccept(orderId uint) bool
	ReservedReject(orderId uint) bool
}
