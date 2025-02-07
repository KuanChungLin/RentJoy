package services

import (
	"fmt"
	"log"
	"math/rand"
	"rentjoy/internal/dto/venuepage"
	repoInterfaces "rentjoy/internal/interfaces/repositories"
	serviceInterfaces "rentjoy/internal/interfaces/services"
	"rentjoy/internal/repositories"

	"gorm.io/gorm"
)

type RecommendedService struct {
	participantRangeRepo repoInterfaces.ParticipantRangeRepository
	venueInformationRepo repoInterfaces.VenueInformationRepository
}

func NewRecommendedService(db *gorm.DB) serviceInterfaces.RecommendedService {
	return &RecommendedService{
		participantRangeRepo: repositories.NewParticipantRangeRepository(db),
		venueInformationRepo: repositories.NewVenueInformationRepository(db),
	}
}

func (s *RecommendedService) GetRecommended() ([]venuepage.Recommended, error) {
	venues, err := s.venueInformationRepo.FindRecommended()
	if err != nil {
		log.Printf("Recommended Get Error: %s", err)
		return []venuepage.Recommended{}, err
	}

	recommended := make([]venuepage.Recommended, 0, 4)
	indexes := rand.Perm(len(venues))

	for i := 0; i < 4 && i < len(indexes); i++ {
		venue := venues[indexes[i]]

		if len(venue.Imgs) == 0 {
			continue
		}

		priceDesc := ""
		if len(venue.BillingRates) > 0 {
			rate := venue.BillingRates[0]
			if rate.RateType.ID == 1 {
				priceDesc = fmt.Sprintf("%s /小時", rate.Rate.StringFixed(0))
			} else {
				priceDesc = fmt.Sprintf("%s /時段", rate.Rate.StringFixed(0))
			}
		}

		recommended = append(recommended, venuepage.Recommended{
			ImgUrl:     venue.Imgs[0].VenueImgPath,
			Name:       venue.Name,
			VenuePrice: priceDesc,
			VenueUrl:   fmt.Sprintf("/Venue/VenuePage?venueId=%d", venue.ID),
		})
	}

	return recommended, nil
}
