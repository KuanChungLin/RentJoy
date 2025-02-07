package interfaces

import "rentjoy/internal/dto/venuepage"

type RecommendedService interface {
	GetRecommended() ([]venuepage.Recommended, error)
}
