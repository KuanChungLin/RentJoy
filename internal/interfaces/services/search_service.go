package interfaces

import (
	"rentjoy/internal/dto/homepage"
	"rentjoy/internal/dto/searchpage"
)

type SearchPageService interface {
	GetSearchPage(searchpage.VenueFilter) *searchpage.SearchPage
	GetActivities() []homepage.Activity
	GetPeopleCounts() []homepage.PeopleCount
	GetSearchPrice() []int
	GetVenueInfos(searchpage.VenueFilter) []searchpage.VenueInfo
}
