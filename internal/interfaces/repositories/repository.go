package interfaces

import (
	"rentjoy/internal/dto/searchpage"
	"rentjoy/internal/models"
	"time"
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
	FindVenuePageByID(venueID int) (*models.VenueInformation, error)
	FindRecommended() ([]models.VenueInformation, error)
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

type DeviceItemRepository interface {
	Repository[models.DeviceItem]
	GetAllDeviceItemNames() ([]string, error)
}

type BillingRateRepository interface {
	Repository[models.BillingRate]
	FindAvailableTimes(venueID int, dayOfWeek time.Weekday) ([]models.BillingRate, error)
	FindByReserved(venueID uint, rateTypeID uint, weekDay int) (*models.BillingRate, error)
}

type OrderRepository interface {
	Repository[models.Order]
	FindConflictingOrders(venueID int, date time.Time) ([]models.Order, error)
}

type VenueImgRepository interface {
	Repository[models.VenueImg]
	FindFirstBySort(venueID uint, sort int) (*models.VenueImg, error)
}
