package services

import (
	"log"
	"rentjoy/internal/dto/venuepage"
	repoInterfaces "rentjoy/internal/interfaces/repositories"
	serviceInterfaces "rentjoy/internal/interfaces/services"
	"rentjoy/internal/repositories"
	"rentjoy/pkg/helper"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PriceService struct {
	billingRateRepo repoInterfaces.BillingRateRepository
}

func NewPriceService(db *gorm.DB) serviceInterfaces.PriceService {
	return &PriceService{
		billingRateRepo: repositories.NewBillingRateRepository(db),
	}
}

// 透過 Id 取得場地時段定價
func (s *PriceService) CalculatePeriodPrice(id int) (int, error) {
	rate, err := s.billingRateRepo.FindByID(uint(id))
	if err != nil {
		log.Printf("Get Rate By Id Error: %s", err)
		return 0, err
	}

	price := helper.DecimalToIntRounded(rate.Rate)

	return price, nil
}

// 透過預訂資料取得價碼
func (s *PriceService) CalculateTimePrices(detail *venuepage.ReservedDetail) (int, error) {
	// 解析日期字串為 time.Time
	reservedDay, err := time.Parse(time.RFC3339, detail.ReservedDay)
	if err != nil {
		log.Printf("Time Parse Error: %s", err)
		return 0, err
	}

	weekDay := int(reservedDay.Weekday())

	rate, err := s.billingRateRepo.FindByReserved(uint(detail.VenueID), 1, weekDay)
	if err != nil {
		log.Printf("Get Rate By Reserved Error: %s", err)
		return 0, nil
	}

	startTime, err := time.Parse(time.RFC3339, detail.StartTime)
	if err != nil {
		return 0, err
	}

	endTime, err := time.Parse(time.RFC3339, detail.EndTime)
	if err != nil {
		return 0, err
	}

	// 檢查是否為 23:59:59
	isEndOfDay := endTime.Hour() == 23 && endTime.Minute() == 59 && endTime.Second() == 59

	var duration time.Duration
	if isEndOfDay {
		// 加一秒計算
		endTimePlusOne := endTime.Add(time.Second)
		duration = endTimePlusOne.Sub(startTime)
	} else {
		duration = endTime.Sub(startTime)
	}

	// 轉換成小時
	hours := decimal.NewFromFloat(duration.Hours())

	// 將 decimal 轉換成 int
	price := helper.DecimalToIntRounded(rate.Rate.Mul(hours))

	return price, nil
}
