package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"rentjoy/internal/dto/order"
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

type OrderService struct {
	participantRangeRepo repoInterfaces.ParticipantRangeRepository
	orderRepo            repoInterfaces.OrderRepository
	orderDetailRepo      repoInterfaces.OrderDetailRepository
	ecpayRepo            repoInterfaces.EcpayRepository
	billingRateRepo      repoInterfaces.BillingRateRepository
	venueRepo            repoInterfaces.VenueInformationRepository
	priceService         serviceInterfaces.PriceService
	ecpayService         serviceInterfaces.EcpayService
	DB                   *gorm.DB
}

func NewOrderService(db *gorm.DB) serviceInterfaces.OrderService {
	return &OrderService{
		participantRangeRepo: repositories.NewParticipantRangeRepository(db),
		orderRepo:            repositories.NewOrderRepository(db),
		orderDetailRepo:      repositories.NewOrderDetailRepository(db),
		ecpayRepo:            repositories.NewEcpayRepository(db),
		billingRateRepo:      repositories.NewBillingRateRepository(db),
		venueRepo:            repositories.NewVenueInformationRepository(db),
		priceService:         NewPriceService(db),
		ecpayService:         NewEcpayService(db),
		DB:                   db,
	}
}

// 產生訂單作業
func (s *OrderService) SaveOrder(orderForm order.OrderForm, userID uint, r *http.Request) (map[string]string, string, error) {
	// 開始交易事務
	tx := s.DB.Begin()
	if tx.Error != nil {
		return nil, "0", tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 資料驗證和準備
	timeDetail, err := s.parseAndValidateTimeDetail(orderForm)
	if err != nil {
		tx.Rollback()
		return nil, "0", err
	}

	// 計算價格
	amount, priceList, err := s.calculatePrice(timeDetail)
	if err != nil {
		tx.Rollback()
		return nil, "0", err
	}

	// 資料型別轉換
	activityTypeID, err := helper.StrToUint(orderForm.Activity)
	if err != nil {
		tx.Rollback()
		return nil, "0", err
	}

	userCount, err := strconv.Atoi(orderForm.UserCount)
	if err != nil {
		tx.Rollback()
		return nil, "0", err
	}

	venue, err := s.venueRepo.FindByID(uint(timeDetail.VenueID))
	if err != nil {
		tx.Rollback()
		return nil, "0", err
	}

	// 建立訂單
	order := &models.Order{
		Status:          0,
		VenueID:         uint(timeDetail.VenueID),
		ActivityTypeID:  activityTypeID,
		FirstName:       orderForm.FirstName,
		LastName:        orderForm.LastName,
		Phone:           orderForm.Phone,
		Email:           orderForm.Email,
		Amount:          amount,
		CreatedAt:       time.Now(),
		UnsubscribeTime: nil,
		UserCount:       userCount,
		Message:         orderForm.Message,
		MemberID:        userID,
		OwnerID:         venue.OwnerID,
	}

	if err := s.orderRepo.CreateOrder(tx, order); err != nil {
		tx.Rollback()
		return nil, "0", err
	}

	// 建立訂單細節
	if err := s.orderDetailRepo.CreateOrderDetails(tx, order.ID, timeDetail, priceList); err != nil {
		tx.Rollback()
		return nil, "0", err
	}

	// 處理 Ecpay 業務
	ecpayParams, err := s.ecpayService.PostOrderDetails(order, r)
	if err != nil {
		tx.Rollback()
		return nil, "0", err
	}

	// 資料型別轉換
	tradeNo, err := strconv.Atoi(ecpayParams["MerchantID"])
	if err != nil {
		tx.Rollback()
		return nil, "0", err
	}

	tradeAmt, err := decimal.NewFromString(ecpayParams["TotalAmount"])
	if err != nil {
		tx.Rollback()
		return nil, "0", err
	}

	// 建立 ECPay 訂單記錄
	ecpayOrder := models.EcpayOrder{
		ID:              order.ID,
		MerchantTradeNo: ecpayParams["MerchantTradeNo"],
		MerchantID:      ecpayParams["MerchantID"],
		RtnCode:         0, // 未付款
		RtnMsg:          "訂單成功尚未付款",
		TradeNo:         tradeNo,
		TradeAmt:        tradeAmt,
		PaymentDate:     time.Now(),
		PaymentType:     ecpayParams["PaymentType"],
		Charge:          decimal.Zero,
		TradeDate:       time.Now(),
		SimulatePaid:    true,
		CheckMacValue:   ecpayParams["CheckMacValue"],
	}

	// 儲存 ECPay 訂單
	if err := s.ecpayRepo.CreateEcpayOrder(tx, ecpayOrder); err != nil {
		tx.Rollback()
		return nil, "0", err
	}

	// 完成交易
	if err := tx.Commit().Error; err != nil {
		return nil, "0", err
	}

	return ecpayParams, strconv.FormatUint(uint64(order.ID), 10), nil
}

// 驗證時間細節
func (s *OrderService) parseAndValidateTimeDetail(form order.OrderForm) (*venuepage.ReservedDetail, error) {
	var detail venuepage.ReservedDetail

	//解析 JSON
	if err := json.Unmarshal([]byte(form.ReservedDetails), &detail); err != nil {
		return nil, fmt.Errorf("parse time detail error: %w", err)
	}

	// 驗證場地 ID 是否一致
	if strconv.Itoa(detail.VenueID) != form.VenueID {
		return nil, errors.New("venue id mismatch")
	}

	// 檢查時間衝突
	isConflict, err := s.isTimeConflict(detail)
	if err != nil {
		return nil, fmt.Errorf("check time conflict error: %w", err)
	}
	if isConflict {
		return nil, errors.New("time slot is already booked")
	}

	return &detail, nil
}

// 確認訂單時間是否衝突
func (s *OrderService) isTimeConflict(timeDetail venuepage.ReservedDetail) (bool, error) {
	// 解析預訂日期
	reservedDay, err := time.Parse(time.RFC3339, timeDetail.ReservedDay)
	if err != nil {
		return false, fmt.Errorf("parse reserved day error: %w", err)
	}

	// 取得已存在的訂單
	existingOrders, err := s.orderRepo.FindConflictingOrders(
		timeDetail.VenueID,
		reservedDay,
	)
	if err != nil {
		return false, fmt.Errorf("find existing orders error: %w", err)
	}

	if timeDetail.StartTime == "" {
		// 時段制預約
		return s.checkPeriodConflict(timeDetail, existingOrders, reservedDay)
	}

	// 小時制預約
	return s.checkHourlyConflict(timeDetail, existingOrders)
}

// 檢查時段制衝突
func (s *OrderService) checkPeriodConflict(timeDetail venuepage.ReservedDetail, existingOrders []models.Order, reservedDay time.Time) (bool, error) {
	// 取得所有要預約的時段資訊
	billingRates, err := s.billingRateRepo.FindByIDs(timeDetail.TimeSlotIds)
	if err != nil {
		return false, err
	}

	// 檢查每個時段是否有衝突
	for _, rate := range billingRates {
		// 預訂日期與預訂時間組合
		bookedTime := helper.CombineDateTime(reservedDay, rate.StartTime, rate.EndTime)

		for _, order := range existingOrders {
			for _, detail := range order.Details {
				if helper.HasTimeOverlap(
					bookedTime.StartTime,
					bookedTime.EndTime,
					detail.StartTime,
					detail.EndTime,
				) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// 檢查小時制衝突
func (s *OrderService) checkHourlyConflict(timeDetail venuepage.ReservedDetail, existingOrders []models.Order) (bool, error) {
	startTime, err := time.Parse(time.RFC3339, timeDetail.StartTime)
	if err != nil {
		return false, err
	}

	endTime, err := time.Parse(time.RFC3339, timeDetail.EndTime)
	if err != nil {
		return false, err
	}

	for _, order := range existingOrders {
		for _, detail := range order.Details {
			if helper.HasTimeOverlap(
				startTime,
				endTime,
				detail.StartTime,
				detail.EndTime,
			) {
				return true, nil
			}
		}
	}

	return false, nil
}

// 計算金額
func (s *OrderService) calculatePrice(timeDetail *venuepage.ReservedDetail) (decimal.Decimal, []decimal.Decimal, error) {
	var totalAmount decimal.Decimal
	var priceList []decimal.Decimal

	if timeDetail.StartTime == "" {
		// 時段制價格計算
		for _, id := range timeDetail.TimeSlotIds {
			intID, err := strconv.Atoi(id)
			if err != nil {
				return decimal.Zero, nil, err
			}

			price, err := s.priceService.CalculatePeriodPrice(intID)
			if err != nil {
				return decimal.Zero, nil, err
			}

			decimalPrice := decimal.NewFromInt(int64(price))
			priceList = append(priceList, decimalPrice)
			totalAmount = totalAmount.Add(decimalPrice)
		}
	} else {
		// 小時制價格計算
		price, err := s.priceService.CalculateTimePrices(timeDetail)
		if err != nil {
			return decimal.Zero, nil, err
		}

		totalAmount = decimal.NewFromInt(int64(price))
		priceList = append(priceList, totalAmount)
	}

	if totalAmount.IsZero() {
		return decimal.Zero, nil, errors.New("invalid price")
	}

	return totalAmount, priceList, nil
}
