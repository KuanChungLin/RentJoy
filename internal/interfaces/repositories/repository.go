package interfaces

import (
	"rentjoy/internal/dto/searchpage"
	"rentjoy/internal/models"
)

type Repository[T any] interface {
	Create(entity T) error
	FindByID(id uint) (*T, error)
	FindAll() ([]T, error)
	Update(entity T) error
	Delete(id uint) error
}

type ParticipantRangeRepository interface {
	Repository[models.ActivityParticipantRange]
}

type ActivityTypeRepository interface {
	Repository[models.ActivityType]
}

type BudgetRepository interface {
	Repository[models.Budget]
}

type VenueInformationRepository interface {
	Repository[models.VenueInformation]
	FindExhibits() ([]models.VenueInformation, error)
	FindExhibitDESC() ([]models.ActivityType, map[uint][]models.VenueInformation, error)
	FindSearchPageInfos(filter searchpage.VenueFilter) ([]models.VenueInformation, error)
}

type MemberRepository interface {
	Repository[models.Member]
	FindByID(id uint) (*models.Member, error)
	FindByAccount(account string) (*models.Member, error)
	IsEmailExists(email string, userID uint) (bool, error)
	FindByFacebookID(ID string) (*models.Member, error)
	FindByGoogleID(ID string) (*models.Member, error)
}

type FacebookRepository interface {
	Repository[models.FacebookThirdPartyLogin]
}

type GoogleRepository interface {
	Repository[models.GoogleThirdPartyLogin]
}
