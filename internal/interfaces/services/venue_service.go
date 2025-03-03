package interfaces

import (
	"rentjoy/internal/dto/venuepage"
	"time"
)

type VenuePageService interface {
	GetVenuePage(venueId int) venuepage.VenuePage
	GetReservedPage(detail *venuepage.ReservedDetail) (venuepage.ReservedPage, error)
	ProcessOrderResult(orderInfo map[string]string) (venuepage.OrderPending, error)
	GetAvailableTime(selectDay time.Time, venueID int) ([]venuepage.AvailableTime, error)
}
