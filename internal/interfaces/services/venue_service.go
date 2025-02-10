package interfaces

import (
	"rentjoy/internal/dto/venuepage"
	"time"
)

type VenuePageService interface {
	GetVenuePage(venueId int) venuepage.VenuePage
	GetReservedPage() venuepage.ReservedPage
	GetOrderPendingPage(orderInfo map[string]string) venuepage.OrderPending
	GetAvailableTime(selectDay time.Time, venueID int) ([]venuepage.AvailableTime, error)
}
