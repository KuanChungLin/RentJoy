package interfaces

import "rentjoy/internal/dto/venuepage"

type VenuePageService interface {
	GetVenuePage(venueId int) venuepage.VenuePage
	GetReservedPage() venuepage.ReservedPage
	GetOrderPendingPage(orderInfo map[string]string) venuepage.OrderPending
}
