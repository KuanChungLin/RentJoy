package services

import (
	"errors"
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

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type VenueService struct {
	participantRangeRepo repoInterfaces.ParticipantRangeRepository
	venueInformationRepo repoInterfaces.VenueInformationRepository
	deviceRepo           repoInterfaces.DeviceItemRepository
	recommendedService   serviceInterfaces.RecommendedService
	billingRateRepo      repoInterfaces.BillingRateRepository
	orderRepo            repoInterfaces.OrderRepository
	ecpayRepo            repoInterfaces.EcpayRepository
	activityRepo         repoInterfaces.ActivityTypeRepository
	venueImgRepo         repoInterfaces.VenueImgRepository
	priceService         serviceInterfaces.PriceService
	ecpayService         serviceInterfaces.EcpayService
	DB                   *gorm.DB
}

func NewVenueService(db *gorm.DB) serviceInterfaces.VenuePageService {
	return &VenueService{
		participantRangeRepo: repositories.NewParticipantRangeRepository(db),
		venueInformationRepo: repositories.NewVenueInformationRepository(db),
		deviceRepo:           repositories.NewDeviceItemRepository(db),
		recommendedService:   NewRecommendedService(db),
		billingRateRepo:      repositories.NewBillingRateRepository(db),
		orderRepo:            repositories.NewOrderRepository(db),
		ecpayRepo:            repositories.NewEcpayRepository(db),
		activityRepo:         repositories.NewActivityTypeRepository(db),
		venueImgRepo:         repositories.NewVenueImgRepository(db),
		priceService:         NewPriceService(db),
		ecpayService:         NewEcpayService(db),
		DB:                   db,
	}
}

// 取得場地資訊頁資料
func (s *VenueService) GetVenuePage(venueID int) venuepage.VenuePage {
	venueChan := make(chan *models.VenueInformation)
	devicesChan := make(chan []string)
	errorChan := make(chan error, 2)

	// 並行查詢場地資訊
	go func() {
		venue, err := s.venueInformationRepo.FindVenuePageByID(uint(venueID))
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
			log.Printf("Goroutine Error: %s", err)
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

	venueInfo := venuepage.VenuePage{
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

	log.Println(venueInfo)
	return venueInfo
}

// 取得預定場地頁資料
func (s *VenueService) GetReservedPage(detail *venuepage.ReservedDetail) (venuepage.ReservedPage, error) {
	// 建立 channels
	venueChan := make(chan *models.VenueInformation)
	imgChan := make(chan *models.VenueImg)
	activitiesChan := make(chan []models.ActivityType)
	errorChan := make(chan error, 3)

	// 並行查詢場地資訊
	go func() {
		venue, err := s.venueInformationRepo.FindByID(uint(detail.VenueID))
		if err != nil {
			errorChan <- err
			return
		}
		venueChan <- venue
	}()

	// 並行查詢圖片
	go func() {
		img, err := s.venueImgRepo.FindFirstBySort(uint(detail.VenueID), 0)
		if err != nil {
			errorChan <- err
			return
		}
		imgChan <- img
	}()

	// 並行查詢活動類型
	go func() {
		activities, err := s.activityRepo.FindAll()
		if err != nil {
			errorChan <- err
			return
		}
		activitiesChan <- activities
	}()

	// 接收結果
	var venueResult *models.VenueInformation
	var imgResult *models.VenueImg
	var activitiesResult []models.ActivityType

	for i := 0; i < 3; i++ {
		select {
		case err := <-errorChan:
			log.Printf("ErrorL %s", err)
		case venueResult = <-venueChan:
		case imgResult = <-imgChan:
		case activitiesResult = <-activitiesChan:
		case <-time.After(30 * time.Second):
			return venuepage.ReservedPage{}, errors.New("timeout getting data")
		}
	}

	if venueResult == nil || imgResult == nil || activitiesResult == nil {
		return venuepage.ReservedPage{}, errors.New("failed to get all required data")
	}

	var timeDetails []venuepage.TimeDetail
	var amount int

	// 處理預訂時間
	if detail.StartTime == "" {
		// 處理時段制
		for _, id := range detail.TimeSlotIds {
			intID, err := strconv.Atoi(id)
			if err != nil {
				return venuepage.ReservedPage{}, err
			}
			rate, err := s.billingRateRepo.FindByID(uint(intID))
			if err != nil {
				return venuepage.ReservedPage{}, err
			}

			price, err := s.priceService.CalculatePeriodPrice(intID)
			if err != nil {
				return venuepage.ReservedPage{}, err
			}

			amount += price

			timeDetails = append(timeDetails, venuepage.TimeDetail{
				TimeRange: fmt.Sprintf("時段 %02d:%02d - %02d:%02d",
					rate.StartTime.Hour(), rate.StartTime.Minute(),
					rate.EndTime.Hour(), rate.EndTime.Minute()),
				Price: strconv.Itoa(price),
			})
		}
	} else {
		// 解析時間
		startTime, _ := time.Parse(time.RFC3339, detail.StartTime)
		endTime, _ := time.Parse(time.RFC3339, detail.EndTime)

		// 處理小時制
		timeRange := ""
		if endTime.Hour() == 23 && endTime.Minute() == 59 {
			timeRange = fmt.Sprintf("小時 %02d:%02d - 24:00",
				startTime.Hour(), startTime.Minute())
		} else {
			timeRange = fmt.Sprintf("小時 %02d:%02d - %02d:%02d",
				startTime.Hour(), startTime.Minute(),
				endTime.Hour(), endTime.Minute())
		}

		price, err := s.priceService.CalculateTimePrices(detail)
		if err != nil {
			return venuepage.ReservedPage{}, err
		}

		amount += price

		timeDetails = append(timeDetails, venuepage.TimeDetail{
			TimeRange: timeRange,
			Price:     strconv.Itoa(price),
		})
	}

	// 解析預訂日期
	reservedDay, _ := time.Parse(time.RFC3339, detail.ReservedDay)
	weekday := helper.GetDayOfWeekInChinese(reservedDay)
	dateStr := fmt.Sprintf("%d 年 %d 月 %d 日 %s",
		reservedDay.Year(), reservedDay.Month(), reservedDay.Day(), weekday)

	return venuepage.ReservedPage{
		VenueID:            strconv.Itoa(detail.VenueID),
		VenueImgUrl:        imgResult.VenueImgPath,
		Name:               venueResult.Name,
		Address:            venueResult.City + venueResult.District + venueResult.Address,
		Date:               dateStr,
		ReservedActivities: helper.ACTModelToDTO(activitiesResult),
		Amount:             strconv.Itoa(amount),
		TimeDetails:        timeDetails,
	}, nil
}

// 取得預訂結果頁資料
func (s *VenueService) GetOrderPendingPage(orderInfo map[string]string) (venuepage.OrderPending, error) {
	// 開始交易
	tx := s.DB.Begin()
	if tx.Error != nil {
		return venuepage.OrderPending{}, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 驗證 CheckMacValue
	ecpayReturnData := s.ecpayService.GetOrderDetails(orderInfo)

	// 取得訂單資訊
	ecpayOrder, err := s.ecpayRepo.FindByMerchantTradeNo(orderInfo["MerchantTradeNo"])
	if err != nil {
		tx.Rollback()
		return venuepage.OrderPending{}, err
	}

	if ecpayOrder == nil {
		tx.Rollback()
		return venuepage.OrderPending{}, errors.New("order not found")
	}

	// 開始更新 Ecpay 訂單
	// 更新 ECPay 訂單資訊
	if err := s.updateEcpayOrder(tx, ecpayOrder, ecpayReturnData); err != nil {
		tx.Rollback()
		log.Printf("Update Ecpay Order Error:%s", err)
		return venuepage.OrderPending{}, err
	}

	// 取得 Order 訂單資料
	order, err := s.orderRepo.FindByEcpayID(tx, ecpayOrder.ID)
	if err != nil {
		tx.Rollback()
		return venuepage.OrderPending{}, err
	}

	// 根據 CheckMacValue 和 RtnCode 返回不同結果
	if ecpayReturnData["CheckMacValue"] == orderInfo["CheckMacValue"] {
		return s.handleValidCheckMacValue(tx, order, orderInfo)
	} else {
		return s.handleInvalidCheckMacValue(tx, order, orderInfo)
	}
}

// 取得場地開放預訂時間
func (s *VenueService) GetAvailableTime(selectDay time.Time, venueID int) ([]venuepage.AvailableTime, error) {
	ratesChan := make(chan []models.BillingRate)
	ordersChan := make(chan []models.Order)
	errorChan := make(chan error, 2)

	// 檢查日期是否有效
	if !helper.IsDateValid(selectDay) {
		return nil, fmt.Errorf("invalid date")
	}

	// 並行查詢計費規則
	go func() {
		rates, err := s.billingRateRepo.FindAvailableTimes(uint(venueID), selectDay.Weekday())
		if err != nil {
			errorChan <- err
			return
		}
		ratesChan <- rates
	}()

	// 並行查詢場地訂單
	go func() {
		orders, err := s.orderRepo.FindConflictingOrders(venueID, selectDay)
		if err != nil {
			errorChan <- err
			return
		}
		ordersChan <- orders
	}()

	var ratesResult []models.BillingRate
	var ordersResult []models.Order

	for i := 0; i < 2; i++ {
		select {
		case err := <-errorChan:
			log.Printf("Goroutine Error: %s", err)
			return nil, err
		case ratesResult = <-ratesChan:
		case ordersResult = <-ordersChan:
		case <-time.After(30 * time.Second):
			return nil, errors.New("timeout getting data")
		}
	}

	// 處理可用時間
	availableTimes := s.processAvailableTime(ratesResult, ordersResult)

	return availableTimes, nil
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

// 更新 EcpayOrder
func (s *VenueService) updateEcpayOrder(tx *gorm.DB, ecpayOrder *models.EcpayOrder, orderInfo map[string]string) error {
	// 更新 EcpayOrder 資料
	ecpayOrder.RtnCode, _ = strconv.Atoi(orderInfo["RtnCode"])
	ecpayOrder.RtnMsg = s.getRtnMsg(orderInfo["RtnMsg"], orderInfo["RtnCode"])
	ecpayOrder.MerchantID = orderInfo["MerchantID"]
	ecpayOrder.TradeAmt, _ = decimal.NewFromString(orderInfo["TradeAmt"])
	ecpayOrder.PaymentDate = s.getPaymentDate(orderInfo["RtnCode"], orderInfo["PaymentDate"])
	ecpayOrder.PaymentType = orderInfo["PaymentType"]
	ecpayOrder.Charge, _ = decimal.NewFromString(orderInfo["PaymentTypeChargeFee"])
	ecpayOrder.TradeDate = helper.MustParseTime(orderInfo["TradeDate"])
	ecpayOrder.SimulatePaid = orderInfo["SimulatePaid"] == "1"
	ecpayOrder.CheckMacValue = orderInfo["CheckMacValue"]

	// 更新到資料庫
	if err := s.ecpayRepo.UpdateByTx(tx, *ecpayOrder); err != nil {
		return fmt.Errorf("update ecpay order error: %w", err)
	}

	return nil
}

// 根據 Ecpay 的 RtnCode 取得 RtnMsg
func (s *VenueService) getRtnMsg(rtnMsg, rtnCode string) string {
	if rtnMsg == "Succeeded" {
		if rtnCode == "1" {
			return "已付款"
		}
		return "付款失敗"
	}
	return rtnMsg
}

// 根據 Ecpay 的 RtnCode 設定 PaymentDate
func (s *VenueService) getPaymentDate(rtnCode, paymentDateStr string) time.Time {
	if rtnCode != "1" {
		return time.Now()
	}
	paymentDate, err := helper.ParseTime(paymentDateStr)
	if err != nil {
		return time.Now()
	}
	return paymentDate
}

// 處理 Ecpay 檢查碼相同時的 Order 資料更新
func (s *VenueService) handleValidCheckMacValue(tx *gorm.DB, order *models.Order, orderInfo map[string]string) (venuepage.OrderPending, error) {
	if orderInfo["RtnCode"] != "1" {
		if err := s.orderRepo.UpdateStatus(tx, order.ID, 5); err != nil {
			return venuepage.OrderPending{}, err
		}
		return venuepage.OrderPending{
			IsPayFail: true,
			VenueId:   strconv.FormatUint(uint64(order.VenueID), 10),
		}, nil
	}

	if err := s.orderRepo.UpdateStatus(tx, order.ID, 1); err != nil {
		return venuepage.OrderPending{}, err
	}
	return venuepage.OrderPending{
		VenueId: strconv.FormatUint(uint64(order.VenueID), 10),
		OrderId: strconv.FormatUint(uint64(order.ID), 10),
		OrderNo: strconv.FormatUint(uint64(order.ID), 10),
		Email:   order.Email,
		Phone:   order.Phone,
	}, nil
}

// 處理 Ecpay 檢查碼不同時的 Order 資料更新
func (s *VenueService) handleInvalidCheckMacValue(tx *gorm.DB, order *models.Order, orderInfo map[string]string) (venuepage.OrderPending, error) {
	if orderInfo["RtnCode"] != "1" {
		if err := s.orderRepo.UpdateStatus(tx, order.ID, 5); err != nil {
			return venuepage.OrderPending{}, err
		}
		return venuepage.OrderPending{
			IsPayFail: true,
			VenueId:   strconv.FormatUint(uint64(order.VenueID), 10),
		}, nil
	}

	if err := s.orderRepo.UpdateStatus(tx, order.ID, 1); err != nil {
		return venuepage.OrderPending{}, err
	}
	return venuepage.OrderPending{
		IsCheckMacValueFail: true,
		VenueId:             strconv.FormatUint(uint64(order.VenueID), 10),
	}, nil
}
