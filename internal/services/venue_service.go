package services

import (
	"rentjoy/internal/dto/venuepage"
	repoInterfaces "rentjoy/internal/interfaces/repositories"
	serviceInterfaces "rentjoy/internal/interfaces/services"
	"rentjoy/internal/repositories"

	"gorm.io/gorm"
)

type VenueService struct {
	participantRangeRepo repoInterfaces.ParticipantRangeRepository
}

func NewVenueService(db *gorm.DB) serviceInterfaces.VenuePageService {
	return &VenueService{
		participantRangeRepo: repositories.NewParticipantRangeRepository(db),
	}
}

func (s *VenueService) GetVenuePage(venueID int) venuepage.VenuePage {
	return venuepage.VenuePage{}
}

func (s *VenueService) GetReservedPage() venuepage.ReservedPage {
	return venuepage.ReservedPage{}
}

func (s *VenueService) GetOrderPendingPage(orderInfo map[string]string) venuepage.OrderPending {
	return venuepage.OrderPending{}
}
