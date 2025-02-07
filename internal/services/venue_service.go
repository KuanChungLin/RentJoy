package services

import (
	"log"
	"rentjoy/internal/dto/venuepage"
	repoInterfaces "rentjoy/internal/interfaces/repositories"
	serviceInterfaces "rentjoy/internal/interfaces/services"
	"rentjoy/internal/models"
	"rentjoy/internal/repositories"
	"rentjoy/pkg/helper"
	"strconv"

	"gorm.io/gorm"
)

type VenueService struct {
	participantRangeRepo repoInterfaces.ParticipantRangeRepository
	venueInformationRepo repoInterfaces.VenueInformationRepository
	deviceRepo           repoInterfaces.DeviceItemRepository
	recommendedService   serviceInterfaces.RecommendedService
}

func NewVenueService(db *gorm.DB) serviceInterfaces.VenuePageService {
	return &VenueService{
		participantRangeRepo: repositories.NewParticipantRangeRepository(db),
		venueInformationRepo: repositories.NewVenueInformationRepository(db),
		deviceRepo:           repositories.NewDeviceItemRepository(db),
		recommendedService:   NewRecommendedService(db),
	}
}

// 取得場地資訊頁資料
func (s *VenueService) GetVenuePage(venueID int) venuepage.VenuePage {
	venueChan := make(chan *models.VenueInformation)
	devicesChan := make(chan []string)
	errorChan := make(chan error, 2)

	// 並行查詢場地資訊
	go func() {
		venue, err := s.venueInformationRepo.FindVenuePageByID(venueID)
		if err != nil {
			errorChan <- err
			return
		}
		venueChan <- venue
	}()

	// 並行查詢設備列表
	go func() {
		devices, err := s.deviceRepo.GetAllDeviceItemNames()
		if err != nil {
			errorChan <- err
			return
		}
		devicesChan <- devices
	}()

	// 接收結果
	var venue *models.VenueInformation
	var devices []string

	// 使用 select 處理兩個 channel
	for i := 0; i < 2; i++ {
		select {
		case err := <-errorChan:
			log.Printf("Error: %s", err)
			return venuepage.VenuePage{}
		case venue = <-venueChan:
		case devices = <-devicesChan:
		}
	}

	recommendedVenues, err := s.recommendedService.GetRecommended()
	if err != nil {
		log.Printf("VenuePage Get recommended Error", err)
		return venuepage.VenuePage{}
	}

	return venuepage.VenuePage{
		VenueID:                 venue.ID,
		ImgUrls:                 helper.GetSortedImgs(venue.Imgs),
		Name:                    venue.Name,
		City:                    venue.City,
		District:                venue.District,
		Address:                 venue.Address,
		NumberOfPeople:          strconv.Itoa(venue.NumOfPeople),
		SpaceSize:               strconv.Itoa(venue.SpaceSize),
		VenueDevices:            helper.ConvertToVenueDevice(venue.Devices),
		VenueNotIncludedDevices: helper.GetNotIncludedDevices(devices, venue.Devices),
		VenueComment:            helper.GetVenueComments(venue.Orders),
		CommentAverage:          venue.EvaluateRate,
		VenueRules:              helper.SplitRules(venue.Rules),
		TrafficInfo: venuepage.TrafficInfo{
			MRTInfo:  venue.MRTInfo,
			BusInfo:  venue.BusInfo,
			ParkInfo: venue.ParkInfo,
		},
		Lng:                venue.Longitude,
		Lat:                venue.Latitude,
		OwnerInfo:          helper.GetOwnerInfo(&venue.Management),
		UnsubscribeRule:    venue.UnsubscribeRule,
		HrPriceRange:       helper.GetPriceRange(1, venue.BillingRates),
		TimeSlotPriceRange: helper.GetPriceRange(2, venue.BillingRates),
		Recommended:        recommendedVenues,
		ReservedDate:       helper.GetReserveDates(venue.Orders),
		OpenDayOfWeek:      helper.GetUniqueDayOfWeek(venue.BillingRates),
		MinRentHours:       helper.GetMinRentHours(venue.BillingRates),
	}
}

func (s *VenueService) GetReservedPage() venuepage.ReservedPage {
	return venuepage.ReservedPage{}
}

func (s *VenueService) GetOrderPendingPage(orderInfo map[string]string) venuepage.OrderPending {
	return venuepage.OrderPending{}
}
