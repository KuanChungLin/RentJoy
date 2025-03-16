package interfaces

import (
	"rentjoy/internal/dto/manage"
)

type ManageService interface {
	GetReservedManagement(userId uint) (*manage.ReservedManagement, error)
	GetVenueManagement(userId uint) (*manage.VenueManagement, error)
	ReservedAccept(orderId uint) bool
	ReservedReject(orderId uint) bool
	DelistVenue(venueId uint) bool
	DeleteVenue(venueId uint) bool
}
