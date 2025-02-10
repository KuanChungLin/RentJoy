package controllers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	interfaces "rentjoy/internal/interfaces/services"
	"rentjoy/pkg/helper"
	"strconv"
	"time"
)

type VenuePageController struct {
	BaseController
	venueService interfaces.VenuePageService
}

func NewVenuePageController(venueService interfaces.VenuePageService, templates map[string]*template.Template) *VenuePageController {
	return &VenuePageController{
		BaseController: NewBaseController(templates),
		venueService:   venueService,
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
		venueID, err = helper.StrToInt(queryID)
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
			venueID, err = helper.StrToInt(formID)
			if err != nil {
				log.Printf("Form 解析錯誤: %s", err)
				// TODO
				// http.Redirect(w, r, "/error", http.StatusSeeOther)
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
