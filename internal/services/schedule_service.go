package services

import (
	"log"
	repoInterfaces "rentjoy/internal/interfaces/repositories"
	serviceInterfaces "rentjoy/internal/interfaces/services"
	"rentjoy/internal/models"
	"rentjoy/internal/repositories"
	"time"

	"gorm.io/gorm"
)

type ScheduleService struct {
	participantRangeRepo repoInterfaces.ParticipantRangeRepository
	DB                   *gorm.DB
}

func NewScheduleService(db *gorm.DB) serviceInterfaces.ScheduleService {
	return &ScheduleService{
		participantRangeRepo: repositories.NewParticipantRangeRepository(db),
		DB:                   db,
	}
}

func (s *ScheduleService) OrderSchedule() {
	// 啟動未付款訂單檢查 (每小時檢查一次)
	go s.runScheduler(time.Hour, func() {
		log.Println("執行未付款訂單檢查...")
		if err := s.expireUnpaidOrders(s.DB); err != nil {
			log.Printf("未付款訂單檢查失敗: %v", err)
		}
	})

	// 啟動未審核訂單檢查
	go s.runScheduler(time.Hour, func() {
		log.Println("執行未審核訂單檢查...")
		if err := s.processUnreviewedOrders(s.DB); err != nil {
			log.Printf("未審核訂單檢查失敗: %v", err)
		}
	})

	// 啟動已完成訂單檢查
	go s.runScheduler(time.Hour, func() {
		log.Println("執行已完成訂單檢查...")
		if err := s.archiveCompletedOrders(s.DB); err != nil {
			log.Printf("已完成訂單檢查失敗: %v", err)
		}
	})

	log.Println("所有排程任務已啟動")
}

// 通用排程執行器
func (s *ScheduleService) runScheduler(interval time.Duration, taskFunc func()) {
	// 立即執行一次
	taskFunc()

	// 創建定時器，取代 Channel 創建
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			taskFunc()
		}
	}
}

// 處理未付款訂單
func (s *ScheduleService) expireUnpaidOrders(db *gorm.DB) error {
	var err error
	// 開始交易事務
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	expirationTime := time.Now().Add(-20 * time.Minute)
	result := db.Model(&models.Order{}).
		Where("OrderStatus = ? AND CreateAt < ?", 0, expirationTime).
		Update("OrderStatus", 5)

	if result.Error != nil {
		return result.Error
	}

	log.Printf("已將 %d 筆未付款訂單標記為過期", result.RowsAffected)

	return nil
}

// 處理未審核訂單
func (s *ScheduleService) processUnreviewedOrders(db *gorm.DB) error {
	var err error
	// 開始交易事務
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	var orderIDs []uint
	overdue := time.Now().Add(24 * time.Hour)

	subQuery := tx.Raw(`
		SELECT o.Id 
		FROM [Order] o 
		JOIN [OrderDetail] od ON o.Id = od.OrderId 
		WHERE o.OrderStatus = 1 
		AND CONVERT(DATE, od.StartTime) <= ? 
		GROUP BY o.Id
	`, overdue.Format("2006-01-02"))

	if err := subQuery.Find(&orderIDs).Error; err != nil {
		return err
	}

	// 沒有找到符合條件的訂單，直接返回
	if len(orderIDs) == 0 {
		return nil
	}

	// 更新這些訂單的狀態和取消時間
	result := tx.Model(&models.Order{}).
		Where("Id IN ?", orderIDs).
		Updates(map[string]interface{}{
			"Status":          3,
			"UnsubscribeTime": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	log.Printf("已更新 %d 筆訂單狀態為已失效", result.RowsAffected)

	return nil
}

// 處理已完成訂單
func (s *ScheduleService) archiveCompletedOrders(db *gorm.DB) error {
	var err error
	// 開始交易事務
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	var completedOrderIDs []uint
	overdue := time.Now()

	subQuery := tx.Raw(`
		SELECT o.Id FROM [Order] o
		WHERE o.OrderStatus = 2
		AND NOT EXISTS (
			SELECT 1 FROM [OrderDetail] od
			WHERE od.OrderId = o.Id
			AND CONVERT(DATE, od.EndTime) >= CONVERT(DATE, ?)
		)
	`, overdue)

	// 將 subQuery 查詢結果映射到 []uint 內
	if err := subQuery.Scan(&completedOrderIDs).Error; err != nil {
		return err
	}

	// 沒有找到符合條件的訂單，直接返回
	if len(completedOrderIDs) == 0 {
		return nil
	}

	// 更新訂單狀態
	result := tx.Model(&models.Order{}).
		Where("Id IN ?", completedOrderIDs).
		Update("Status", 4)

	if result.Error != nil {
		return result.Error
	}

	log.Printf("已更新 %d 筆訂單狀態為已結束", result.RowsAffected)

	return nil
}
