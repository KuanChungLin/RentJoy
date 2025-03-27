package controllers

import (
	"html/template"
	"log"
	"net/http"
	interfaces "rentjoy/internal/interfaces/services"
	"rentjoy/internal/middleware"
	"rentjoy/pkg/helper"
)

type ManageController struct {
	BaseController
	manageService interfaces.ManageService
}

func NewManageController(manageService interfaces.ManageService, templates map[string]*template.Template) *ManageController {
	return &ManageController{
		BaseController: NewBaseController(templates),
		manageService:  manageService,
	}
}

// 預訂單管理頁面
func (c *ManageController) ReservedManagement(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("Get userId Error")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vm, err := c.manageService.GetReservedManagement(userId)
	if err != nil {
		http.Error(w, "Reserved Information Get Error", http.StatusNotImplemented)
	}

	c.RenderTemplate(w, r, "manage_reserved", vm)
}

// 場地管理頁面
func (c *ManageController) VenueManagement(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vm, err := c.manageService.GetVenueManagement(userId)
	if err != nil {
		http.Error(w, "Venue Managemant Get Error", http.StatusNotImplemented)
	}

	c.RenderTemplate(w, r, "manage_venues", vm)
}

// 預訂單接受預訂作業
func (c *ManageController) ReservedAccept(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "無法解析表單數據", http.StatusBadRequest)
		return
	}

	orderId, err := helper.StrToUint(r.FormValue("orderId"))
	if err != nil {
		log.Printf("Get FormValue OrderId Error:%s", err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Fail"))
		return
	}

	ok := c.manageService.ReservedAccept(orderId)
	if !ok {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Fail"))
		return
	}

	// 設置響應頭為純文本
	w.Header().Set("Content-Type", "text/plain")
	// 設置狀態碼為 200 OK
	w.WriteHeader(http.StatusOK)
	// 寫入響應內容
	w.Write([]byte("Success"))
}

// 預訂單拒絕預訂作業
func (c *ManageController) ReservedReject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "無法解析表單數據", http.StatusBadRequest)
		return
	}

	orderId, err := helper.StrToUint(r.FormValue("orderId"))
	if err != nil {
		log.Printf("Get FormValue OrderId Error:%s", err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Fail"))
		return
	}

	ok := c.manageService.ReservedReject(orderId)
	if !ok {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Fail"))
		return
	}

	// 設置響應頭為純文本
	w.Header().Set("Content-Type", "text/plain")
	// 設置狀態碼為 200 OK
	w.WriteHeader(http.StatusOK)
	// 寫入響應內容
	w.Write([]byte("Success"))
}

// 場地下架作業
func (c *ManageController) DelistVenue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "無法解析表單數據", http.StatusBadRequest)
		return
	}

	venueId, err := helper.StrToUint(r.FormValue("venueId"))
	if err != nil {
		http.Error(w, "無法解析 venueId", http.StatusBadRequest)
		return
	}

	ok := c.manageService.DelistVenue(venueId)
	if !ok {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Fail"))
		return
	}

	// 設置響應頭為純文本
	w.Header().Set("Content-Type", "text/plain")
	// 設置狀態碼為 200 OK
	w.WriteHeader(http.StatusOK)
	// 寫入響應內容
	w.Write([]byte("Success"))
}

// 場地刪除作業
func (c *ManageController) DeleteVenue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "無法解析表單數據", http.StatusBadRequest)
		return
	}

	venueId, err := helper.StrToUint(r.FormValue("venueId"))
	if err != nil {
		http.Error(w, "無法解析 venueId", http.StatusBadRequest)
		return
	}

	ok := c.manageService.DeleteVenue(venueId)
	if !ok {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Fail"))
		return
	}

	// 設置響應頭為純文本
	w.Header().Set("Content-Type", "text/plain")
	// 設置狀態碼為 200 OK
	w.WriteHeader(http.StatusOK)
	// 寫入響應內容
	w.Write([]byte("Success"))
}
