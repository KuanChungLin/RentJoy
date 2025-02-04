package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	// 讀取 .env 檔案
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file: %w", err)
	}

	// 設定資料庫連線字串
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&encrypt=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	// 建立資料庫連線
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 獲得底層的 SQL DB 物件
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	// 設定空閒連線池的最大連線數
	sqlDB.SetMaxIdleConns(10)
	// 設定資料庫的最大開放連線數
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}
