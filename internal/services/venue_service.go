package services

import (
	"fmt"
	"log"
	"rentjoy/internal/dto/venuepage"
	repoInterfaces "rentjoy/internal/interfaces/repositories"
	serviceInterfaces "rentjoy/internal/interfaces/services"
	"rentjoy/internal/models"
	"rentjoy/internal/repositories"
	"rentjoy/pkg/helper"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type VenueService struct {
	participantRangeRepo repoInterfaces.ParticipantRangeRepository
	venueInformationRepo repoInterfaces.VenueInformationRepository
	deviceRepo           repoInterfaces.DeviceItemRepository
	recommendedService   serviceInterfaces.RecommendedService
	billingRateRepo      repoInterfaces.BillingRateRepository
	orderRepo            repoInterfaces.OrderRepository
}

func NewVenueService(db *gorm.DB) serviceInterfaces.VenuePageService {
	return &VenueService{
		participantRangeRepo: repositories.NewParticipantRangeRepository(db),
		venueInformationRepo: repositories.NewVenueInformationRepository(db),
		deviceRepo:           repositories.NewDeviceItemRepository(db),
		recommendedService:   NewRecommendedService(db),
		billingRateRepo:      repositories.NewBillingRateRepository(db),
		orderRepo:            repositories.NewOrderRepository(db),
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
		log.Printf("VenuePage Get recommended Error: %s", err)
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

// 取得場地開放預訂時間
func (s *VenueService) GetAvailableTime(selectDay time.Time, venueID int) ([]venuepage.AvailableTime, error) {
	// 檢查日期是否有效
	if !helper.IsDateValid(selectDay) {
		return nil, fmt.Errorf("invalid date")
	}

	// 取得計費規則
	rates, err := s.billingRateRepo.FindAvailableTimes(venueID, selectDay.Weekday())
	if err != nil {
		log.Printf("AvailableTimes Get Error: %s", err)
		return nil, err
	}

	// 獲取場地訂單
	orders, err := s.orderRepo.FindConflictingOrders(venueID, selectDay)
	if err != nil {
		log.Printf("ConflictingOrders Get Error: %s", err)
		return nil, err
	}

	// 處理可用時間
	availableTimes := s.processAvailableTime(rates, orders)

	return availableTimes, err
}

// 處理可用時間
func (s *VenueService) processAvailableTime(rates []models.BillingRate, orders []models.Order) []venuepage.AvailableTime {
	var availableTimes []venuepage.AvailableTime

	for _, rate := range rates {
		switch rate.RateTypeID {
		case 1: // 小時制
			times := s.processHourlyRate(rate, orders)
			availableTimes = append(availableTimes, times...)
		case 2: // 時段制
			if time := s.processTimeSlotRate(rate, orders); time != nil {
				availableTimes = append(availableTimes, *time)
			}
		}
	}

	return availableTimes
}

// 處理小時制時間
func (s *VenueService) processHourlyRate(rate models.BillingRate, orders []models.Order) []venuepage.AvailableTime {
	var times []venuepage.AvailableTime
	currentTime := rate.StartTime

	for currentTime.Before(rate.EndTime) {
		endTime := currentTime.Add(30 * time.Minute)
		if endTime.After(rate.EndTime) {
			endTime = rate.EndTime
		}

		if !helper.IsTimeConflict(currentTime, endTime, orders) {
			times = append(times, venuepage.AvailableTime{
				StartTime:     currentTime.Format("15:04"),
				EndTime:       endTime.Format("15:04"),
				Price:         rate.Rate.StringFixed(0),
				BillingRateID: strconv.FormatUint(uint64(rate.ID), 10),
				RateTypeID:    strconv.FormatUint(uint64(rate.RateTypeID), 10),
			})
		}

		currentTime = endTime
	}

	return times
}

// 處理時段制時間
func (s *VenueService) processTimeSlotRate(rate models.BillingRate, orders []models.Order) *venuepage.AvailableTime {
	// 檢查此時段是否有衝突
	for _, order := range orders {
		for _, detail := range order.Details {
			// 將時間統一到同一天比較
			detailStart := helper.NormalizeTime(detail.StartTime)
			detailEnd := helper.NormalizeTime(detail.EndTime)
			rateStart := helper.NormalizeTime(rate.StartTime)
			rateEnd := helper.NormalizeTime(rate.EndTime)

			// 時段制要求完全相同才算衝突
			if detailStart.Equal(rateStart) && detailEnd.Equal(rateEnd) {
				return nil
			}
		}
	}

	// 無衝突，返回可用時段
	return &venuepage.AvailableTime{
		StartTime:     rate.StartTime.Format("15:04"),
		EndTime:       rate.EndTime.Format("15:04"),
		Price:         rate.Rate.StringFixed(0),
		BillingRateID: strconv.FormatUint(uint64(rate.ID), 10),
		RateTypeID:    strconv.FormatUint(uint64(rate.RateTypeID), 10),
	}
}
