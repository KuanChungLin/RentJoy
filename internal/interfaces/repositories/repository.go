package interfaces

import (
	"rentjoy/internal/dto/order"
	"rentjoy/internal/dto/searchpage"
	"rentjoy/internal/dto/venuepage"
	"rentjoy/internal/models"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
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
	FindVenuePageByID(venueID uint) (*models.VenueInformation, error)
	FindRecommended() ([]models.VenueInformation, error)
	FindByOwnerId(ownerId uint) ([]models.VenueInformation, error)
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
	FindAvailableTimes(venueID uint, dayOfWeek time.Weekday) ([]models.BillingRate, error)
	FindByReserved(venueID uint, rateTypeID uint, weekDay int) (*models.BillingRate, error)
	FindByIDs(ids []string) ([]models.BillingRate, error)
}

type OrderRepository interface {
	Repository[models.Order]
	CreateOrder(tx *gorm.DB, order *models.Order) error
	FindConflictingOrders(venueID int, date time.Time) ([]models.Order, error)
	FindByUserAndStatus(userId uint, status order.OrderStatus, pageIndex int, pageSize int) ([]models.Order, error)
	FindByEcpayID(tx *gorm.DB, id uint) (*models.Order, error)
	CountByUserAndStatus(userId uint, status order.OrderStatus) (int, error)
	UpdateStatus(tx *gorm.DB, id uint, status int) error
	FindOrdersByVenueId(tx *gorm.DB, venueId uint) ([]models.Order, error)
	FindManageOrderByUserId(userId uint) ([]models.Order, error)
}

type OrderDetailRepository interface {
	Repository[models.OrderDetail]
	CreateOrderDetails(tx *gorm.DB, orderID uint, timeDetail *venuepage.ReservedDetail, priceList []decimal.Decimal) error
	FindByOrderID(orderId uint) ([]models.OrderDetail, error)
}

type EcpayRepository interface {
	Repository[models.EcpayOrder]
	CreateEcpayOrder(tx *gorm.DB, ecpayOrder models.EcpayOrder) error
	FindByMerchantTradeNo(tradeNo string) (*models.EcpayOrder, error)
	UpdateByTx(tx *gorm.DB, ecpayOrder models.EcpayOrder) error
}

type VenueImgRepository interface {
	Repository[models.VenueImg]
	FindFirstBySort(venueID uint, sort int) (*models.VenueImg, error)
}

type VenueEvaluateRepository interface {
	Repository[models.VenueEvaluate]
	FindByOrderId(orderId uint) (*models.VenueEvaluate, error)
	CreateByTx(tx *gorm.DB, evaluate models.VenueEvaluate) error
}

type SpaceTypeRepository interface {
	Repository[models.SpaceType]
	FindSpaceAndFacility() ([]models.SpaceType, error)
}

type DeviceTypeRepository interface {
	Repository[models.DeviceType]
	FindTypeAndItems() ([]models.DeviceType, error)
}

type ManagementRepository interface {
	Repository[models.Management]
	FindByUserId(userId uint) ([]models.Management, error)
}
