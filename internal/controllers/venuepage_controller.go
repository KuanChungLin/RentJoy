package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"rentjoy/internal/dto/venuepage"
	interfaces "rentjoy/internal/interfaces/services"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type VenuePageController struct {
	BaseController
	venueService interfaces.VenuePageService
	redisClient  *redis.Client
}

func NewVenuePageController(venueService interfaces.VenuePageService, templates map[string]*template.Template, redisClient *redis.Client) *VenuePageController {
	return &VenuePageController{
		BaseController: NewBaseController(templates),
		venueService:   venueService,
		redisClient:    redisClient,
	}
}

// 場地資訊頁
func (c *VenuePageController) VenuePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var venueID int
	var err error

	// 檢查 Query String
	if queryID := r.URL.Query().Get("venueId"); queryID != "" {
		venueID, err = strconv.Atoi(queryID)
		if err != nil {
			log.Printf("Query String 解析錯誤: %s", err)
			http.Redirect(w, r, "/error", http.StatusSeeOther)
			return
		}
	} else {
		// 檢查 Form
		if err := r.ParseForm(); err != nil {
			http.Error(w, "無法解析表單數據", http.StatusBadRequest)
			return
		}

		if formID := r.FormValue("venueId"); formID != "" {
			venueID, err = strconv.Atoi(formID)
			if err != nil {
				log.Printf("Form 解析錯誤: %s", err)
				http.Redirect(w, r, "/error", http.StatusSeeOther)
				return
			}
		}
	}

	// 驗證 venueID
	if venueID <= 0 {
		log.Printf("無效的 venueID: %d", venueID)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	// 刪除 Cookie
	cookie := &http.Cookie{
		Name:   "TimeDetailCookie",
		Path:   "/Venue",
		MaxAge: -1, // 設為-1表示立即刪除
	}
	http.SetCookie(w, cookie)

	vm := c.venueService.GetVenuePage(venueID)

	c.RenderTemplate(w, r, "venuepage", vm)
}

// 取得場地空缺時間
func (c *VenuePageController) GetAvailableTime(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	selectDayStr := r.URL.Query().Get("selectDay")
	venueIDStr := r.URL.Query().Get("venueId")

	// 解析時間
	selectDay, err := time.Parse(time.RFC3339, selectDayStr)
	if err != nil {
		log.Printf("日期解析錯誤: %v", err)
		http.Error(w, "無法解析日期", http.StatusBadRequest)
		return
	}

	// 解析 venueId
	venueID, err := strconv.Atoi(venueIDStr)
	if err != nil {
		log.Printf("ID解析錯誤: %v", err)
		http.Error(w, "無效的場地ID", http.StatusBadRequest)
		return
	}

	data, err := c.venueService.GetAvailableTime(selectDay, venueID)
	if err != nil {
		log.Printf("取得開放預約時間錯誤: %s", err)
		http.Error(w, "服務器錯誤", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON編碼錯誤: %v", err)
		http.Error(w, "服務器錯誤", http.StatusInternalServerError)
		return
	}
}

// 預訂場地頁面
func (c *VenuePageController) ReservedPage(w http.ResponseWriter, r *http.Request) {
	// 檢查 Ｃookie 是否存在
	cookie, err := r.Cookie("TimeDetailCookie")
	if err == http.ErrNoCookie {
		log.Printf("TimeDetailCookie Not Exist Error: %s", err)
		http.Redirect(w, r, "/Error", http.StatusSeeOther)
		return
	}

	// URL 解碼
	decodedValue, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		log.Printf("Cookie decode error: %s", err)
		// http.Redirect(w, r, "/Error", http.StatusSeeOther)
		return
	}

	// 解析取得 Cookie 資料
	var timeDetail venuepage.ReservedDetail
	if err := json.Unmarshal([]byte(decodedValue), &timeDetail); err != nil {
		log.Printf("JSON Unmarshal Error: %s", err)
		// http.Redirect(w, r, "/Error", http.StatusSeeOther)
		return
	}

	vm, err := c.venueService.GetReservedPage(&timeDetail)
	if err != nil {
		http.Redirect(w, r, "/Error", http.StatusSeeOther)
		return
	}

	vm.ReservedDetailCookie = decodedValue

	// 刪除 Cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "TimeDetailCookie",
		MaxAge: -1,
	})

	c.RenderTemplate(w, r, "reservedpage", vm)
}

// 預訂場地結果頁
func (c *VenuePageController) OrderPending(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 從 QueryString 取得訂單 ID
	tradeNo := r.URL.Query().Get("merchantTradeNo")
	if tradeNo == "" {
		http.Error(w, "Invalid merchantTradeNo", http.StatusBadRequest)
		return
	}

	// 從 Redis 取得訂單資訊
	key := fmt.Sprintf("order:%s", tradeNo)
	orderData, err := c.redisClient.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			http.Error(w, "Order not found", http.StatusNotFound)
		} else {
			log.Printf("Redis Get Error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// 解析 JSON 字串回 map
	var orderInfo map[string]string
	if err := json.Unmarshal([]byte(orderData), &orderInfo); err != nil {
		log.Printf("JSON Unmarshal Error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 如果存在編碼過的 CheckMacValue，則解碼恢復
	if encodedCheckMac, exists := orderInfo["CheckMacValue"]; exists {
		decodedBytes, err := base64.StdEncoding.DecodeString(encodedCheckMac)
		if err == nil {
			// 將解碼後的值寫回原始鍵
			orderInfo["CheckMacValue"] = string(decodedBytes)
		} else {
			log.Printf("Base64 Decode Error: %v", err)
		}
	}

	// 透過 Service 處理訂單資訊
	pendingOrder, err := c.venueService.ProcessOrderResult(orderInfo)
	if err != nil {
		log.Printf("Get Order Pending Info Error: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	c.RenderTemplate(w, r, "order_pending", pendingOrder)
}
