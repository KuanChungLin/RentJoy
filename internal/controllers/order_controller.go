package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"rentjoy/internal/dto/order"
	"rentjoy/internal/dto/venuepage"
	interfaces "rentjoy/internal/interfaces/services"
	"rentjoy/internal/middleware"
	"rentjoy/pkg/helper"
	"time"

	"github.com/go-redis/redis"
)

type OrderController struct {
	BaseController
	orderService interfaces.OrderService
	redisClient  *redis.Client
}

func NewOrderController(orderService interfaces.OrderService, templates map[string]*template.Template, redisClient *redis.Client) *OrderController {
	return &OrderController{
		BaseController: NewBaseController(templates),
		orderService:   orderService,
		redisClient:    redisClient,
	}
}

// 創建訂單記錄並跳轉 Ecpay
func (c *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析表單
	if err := r.ParseForm(); err != nil {
		http.Error(w, "無法解析表單數據", http.StatusBadRequest)
		return
	}

	orderRequest := order.OrderForm{
		Activity:        r.FormValue("Activity"),
		UserCount:       r.FormValue("UserCount"),
		Message:         r.FormValue("Message"),
		LastName:        r.FormValue("LastName"),
		FirstName:       r.FormValue("FirstName"),
		Phone:           r.FormValue("Phone"),
		Email:           r.FormValue("Email"),
		VenueID:         r.FormValue("VenueId"),
		ReservedDetails: r.FormValue("ReservedDetails"),
	}

	// 取得使用者 ID
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 從 Service 獲取 ECPay 參數和訂單 ID
	orderData, orderID, err := c.orderService.SaveOrder(orderRequest, userID, r)
	if err != nil || orderData == nil {
		log.Printf("Save Order Error:%s", err)

		failOrder := venuepage.OrderPending{
			IsCreateFail: true,
			VenueId:      orderRequest.VenueID,
		}
		c.RenderTemplate(w, r, "order_pending", failOrder)
		return
	}

	// 存儲前進行 CheckMacValue 的 Base64 編碼
	if checkMacValue, exists := orderData["CheckMacValue"]; exists {
		// 將原始 CheckMacValue 進行 Base64 編碼
		encodedCheckMac := base64.StdEncoding.EncodeToString([]byte(checkMacValue))
		// 保存原始值到一個特殊的鍵中
		orderData["CheckMacValue"] = encodedCheckMac
	}

	// 將 map 轉換為 JSON
	dataJSON, err := json.Marshal(orderData)
	if err != nil {
		log.Printf("JSON Marshal Error: %v", err)
		return
	}

	key := fmt.Sprintf("order:%s", orderID)

	// 存入 Redis，15分鐘過期
	err = c.redisClient.Set(key, dataJSON, 15*time.Minute).Err()
	if err != nil {
		log.Printf("Redis Set Error: %v", err)
		failOrder := venuepage.OrderPending{
			IsCreateFail: true,
			VenueId:      orderRequest.VenueID,
		}
		c.RenderTemplate(w, r, "order_pending", failOrder)
		return
	}

	http.Redirect(w, r, "/Ecpay/Process?id="+orderID, http.StatusSeeOther)
}

// 訂單處理中資料顯示
func (c *OrderController) OrderReserved(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("Get userId Error")
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}
	pageIndex := 1
	pageSize := 5

	orderInfo, err := c.orderService.GetOrderPage(userId, order.Reserved, pageIndex, pageSize)
	if err != nil {
		log.Printf("Get OrderPage Error:%s", err)
	}

	c.RenderTemplate(w, r, "order_reserved", orderInfo)
}

// 訂單已預訂資料顯示
func (c *OrderController) OrderProcessing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("Get userId Error")
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}
	pageIndex := 1
	pageSize := 5

	orderInfo, err := c.orderService.GetOrderPage(userId, order.Processing, pageIndex, pageSize)
	if err != nil {
		log.Printf("Get OrderPage Error:%s", err)
		return
	}

	c.RenderTemplate(w, r, "order_reserved", orderInfo)
}

// 訂單退訂資料顯示
func (c *OrderController) OrderCancel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("Get userId Error")
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}
	pageIndex := 1
	pageSize := 5

	orderInfo, err := c.orderService.GetOrderPage(userId, order.Cancel, pageIndex, pageSize)
	if err != nil {
		log.Printf("Get OrderPage Error:%s", err)
		return
	}

	c.RenderTemplate(w, r, "order_reserved", orderInfo)
}

// 訂單已結束資料顯示
func (c *OrderController) OrderFinished(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("Get userId Error")
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}
	pageIndex := 1
	pageSize := 5

	orderInfo, err := c.orderService.GetOrderPage(userId, order.Finished, pageIndex, pageSize)
	if err != nil {
		log.Printf("Get OrderPage Error:%s", err)
		return
	}

	c.RenderTemplate(w, r, "order_reserved", orderInfo)
}

// 取消預訂
func (c *OrderController) CancelReservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析表單
	if err := r.ParseForm(); err != nil {
		http.Error(w, "無法解析表單數據", http.StatusBadRequest)
		return
	}

	orderId, err := helper.StrToUint(r.FormValue("orderId"))
	if err != nil {
		log.Printf("無法解析 OrderId Error:%s", err)
		return
	}

	err = c.orderService.CancelReservation(orderId)
	if err != nil {
		log.Printf("Cancel Reservation Error:%s", err)
		http.Redirect(w, r, "/Order/OrderReserved", http.StatusSeeOther)
	}

	http.Redirect(w, r, "/Order/OrderReserved", http.StatusSeeOther)
}
