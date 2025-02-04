package interfaces

import (
	"rentjoy/internal/dto/homepage"
)

type HomeService interface {
	GetHomePage() homepage.HomePage
	GetActivities() []homepage.Activity
	GetPeopleCounts() []homepage.PeopleCount
	GetGalleries() []homepage.Gallery
	GetExhibits() []homepage.Exhibit
	GetExhibitDESC() []homepage.ExhibitDESC
}
