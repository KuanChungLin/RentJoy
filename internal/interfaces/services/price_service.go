package interfaces

import "rentjoy/internal/dto/venuepage"

type PriceService interface {
	CalculatePeriodPrice(id int) (int, error)
	CalculateTimePrices(*venuepage.ReservedDetail) (int, error)
}
