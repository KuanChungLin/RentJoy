package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
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
	orderRepo         repoInterfaces.OrderRepository
	orderDetailRepo   repoInterfaces.OrderDetailRepository
	ecpayRepo         repoInterfaces.EcpayRepository
	billingRateRepo   repoInterfaces.BillingRateRepository
	venueRepo         repoInterfaces.VenueInformationRepository
	venueImgRepo      repoInterfaces.VenueImgRepository
	venueEvaluateRepo repoInterfaces.VenueEvaluateRepository
	priceService      serviceInterfaces.PriceService
	ecpayService      serviceInterfaces.EcpayService
	recommendService  serviceInterfaces.RecommendedService
	DB                *gorm.DB
}

func NewOrderService(db *gorm.DB) serviceInterfaces.OrderService {
	return &OrderService{
		orderRepo:         repositories.NewOrderRepository(db),
		orderDetailRepo:   repositories.NewOrderDetailRepository(db),
		ecpayRepo:         repositories.NewEcpayRepository(db),
		billingRateRepo:   repositories.NewBillingRateRepository(db),
		venueRepo:         repositories.NewVenueInformationRepository(db),
		venueImgRepo:      repositories.NewVenueImgRepository(db),
		venueEvaluateRepo: repositories.NewVenueEvaluateRepository(db),
		priceService:      NewPriceService(db),
		ecpayService:      NewEcpayService(db),
		recommendService:  NewRecommendedService(db),
		DB:                db,
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

// 根據訂單狀態取得預訂場地相關資訊
// orderStatus 訂單狀態
// pageIndex   當前頁數
// pageSize    每頁顯示訂單筆數
func (s *OrderService) GetOrderPage(userId uint, orderStatus order.OrderStatus, pageIndex int, pageSize int) (order.OrderPageInfo, error) {
	// 取得該使用者的所有指定狀態下的預訂次數
	orderCount, err := s.orderRepo.CountByUserAndStatus(userId, orderStatus)
	if err != nil {
		log.Printf("Get Orders Error:%s", err)
		return order.OrderPageInfo{}, err
	}

	// 取得該使用者的所有指定狀態下的預訂記錄
	orders, err := s.orderRepo.FindByUserAndStatus(userId, orderStatus, pageIndex, pageSize)
	if err != nil {
		log.Printf("Get Orders Error:%s", err)
		return order.OrderPageInfo{}, err
	}

	// 根據預訂紀錄取得介面要顯示的相關資訊
	vm, err := s.convertToOrderVM(orders)
	if err != nil {
		log.Printf("Get OrderVM Error:%s", err)
		return order.OrderPageInfo{}, err
	}

	// 取得推薦場地資訊
	recommendedVenues, err := s.recommendService.GetRecommended()
	if err != nil {
		log.Printf("Get Recommended Venues Error:%s", err)
		return order.OrderPageInfo{}, err
	}

	totalPages := int(math.Ceil(float64(orderCount) / float64(pageSize)))

	return order.OrderPageInfo{
		Orders:      vm,
		Recommend:   recommendedVenues,
		OrderCount:  orderCount,
		TotalPages:  totalPages,
		CurrentPage: pageIndex,
	}, nil
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

// 取得預訂單資訊細節
func (s *OrderService) convertToOrderVM(orders []models.Order) ([]order.OrderVM, error) {
	var orderVMs []order.OrderVM

	for _, o := range orders {
		// 獲取場地信息
		venue, err := s.venueRepo.FindByID(o.VenueID)
		if err != nil {
			log.Printf("Venue not found for order %d: %s", o.ID, err)
			continue
		}

		// 獲取場地圖片
		venueImg, err := s.venueImgRepo.FindFirstBySort(o.VenueID, 0)
		if err != nil {
			log.Printf("Venue image not found for venue %d: %s", o.VenueID, err)
			continue
		}

		// 獲取評價
		var orderEvaluate order.OrderEvaluate
		hasEvaluate := false

		// 設定默認值
		orderEvaluate = order.OrderEvaluate{
			Rating:       0,  // 默認評分為0
			Content:      "", // 默認評價內容為空字符串
			EvaluateTime: "", // 默認評價時間為空字符串
		}

		evaluate, err := s.venueEvaluateRepo.FindByOrderId(o.ID)
		if err == nil && evaluate != nil {
			orderEvaluate = order.OrderEvaluate{
				Rating:       int(evaluate.EvaluateRate),
				Content:      evaluate.EvaluateComment,
				EvaluateTime: evaluate.CreatedAt.Format("2006-01-02 15:04:05"),
			}
			hasEvaluate = true
		}

		// 獲取預訂時間
		orderDetails, err := s.orderDetailRepo.FindByOrderID(o.ID)
		if err != nil {
			log.Printf("Order details not found for order %d: %s", o.ID, err)
			continue
		}

		// 創建預訂時間列表
		var scheduledTimes []order.OrderScheduleTime
		for _, detail := range orderDetails {
			scheduledTimes = append(scheduledTimes, order.OrderScheduleTime{
				StartTime: detail.StartTime.Format("15:04"),
				EndTime:   detail.EndTime.Format("15:04"),
			})
		}

		// 取得退訂時間
		var unsubscribeTime string
		if o.UnsubscribeTime != nil {
			unsubscribeTime = o.UnsubscribeTime.Format("2006-01-02 15:04")
		}

		// 創建訂單視圖模型
		orderVM := order.OrderVM{
			Title:          venue.Name,
			Address:        venue.Address,
			OrderId:        fmt.Sprintf("%d", o.ID),
			OrderTime:      o.CreatedAt.Format("2006-01-02 15:04:05"),
			CancelTime:     unsubscribeTime,
			Status:         order.OrderStatus(o.Status),
			OrderPrice:     helper.DecimalToIntRounded(o.Amount),
			ContactPerson:  o.LastName + o.FirstName,
			Email:          o.Email,
			ImgUrl:         venueImg.VenueImgPath,
			ProductUrl:     fmt.Sprintf("/Venue/VenuePage?venueId=%d", o.VenueID),
			ScheduledTimes: scheduledTimes,
		}

		// 只有當評價存在時才設置評價
		if hasEvaluate {
			orderVM.Evaluate = orderEvaluate
		}

		orderVMs = append(orderVMs, orderVM)
	}

	return orderVMs, nil
}
