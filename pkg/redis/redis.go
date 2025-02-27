package redis

import (
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis"
)

func NewRedisClient() *redis.Client {
	// 取得環境變數，設定預設值
	addr := getEnvOrDefault("REDIS_ADDR", "localhost:6379")
	password := getEnvOrDefault("REDIS_PASSWORD", "")
	dbStr := getEnvOrDefault("REDIS_DB", "0")

	// 轉換 DB 編號
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		log.Printf("Invalid REDIS_DB value, using default 0: %v", err)
		db = 0
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 測試連線
	_, err = client.Ping().Result()
	if err != nil {
		log.Printf("Redis connection error: %v", err)
	}

	return client
}

// 從 .env檔 取得設定值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
