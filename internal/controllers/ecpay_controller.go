package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"rentjoy/internal/dto/venuepage"
	"time"

	"github.com/go-redis/redis"
)

type EcpayController struct {
	BaseController
	redisClient *redis.Client
}

func NewEcpayController(templates map[string]*template.Template, redisClient *redis.Client) *EcpayController {
	return &EcpayController{
		BaseController: NewBaseController(templates),
		redisClient:    redisClient,
	}
}

// 跳轉 Ecpay 頁面處理
func (c *EcpayController) Process(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orderID := r.URL.Query().Get("id")
	if orderID == "" {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	key := fmt.Sprintf("order:%s", orderID)

	// 從 Redis 取得資料
	ecpayParams, err := c.redisClient.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			// 找不到資料
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		log.Printf("Redis Get Error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 刪除已取得的資料
	// c.redisClient.Del(key)

	// 解析 JSON 字串回 map
	var orderData map[string]string
	if err := json.Unmarshal([]byte(ecpayParams), &orderData); err != nil {
		log.Printf("JSON Unmarshal Error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 如果存在編碼過的 CheckMacValue，則解碼恢復
	if encodedCheckMac, exists := orderData["CheckMacValue"]; exists {
		decodedBytes, err := base64.StdEncoding.DecodeString(encodedCheckMac)
		if err == nil {
			// 將解碼後的值寫回原始鍵
			orderData["CheckMacValue"] = string(decodedBytes)
		} else {
			log.Printf("Base64 Decode Error: %v", err)
		}
	}

	c.RenderTemplate(w, r, "ecpay_process", orderData)
}

// Ecpay 回傳付款結果處理
func (c *EcpayController) ReceivePaymentResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析表單資料
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// 取得重要資訊儲存到 Redis 的資料
	orderInfo := map[string]string{
		"CustomField1":         r.FormValue("CustomField1"),
		"CustomField2":         r.FormValue("CustomField2"),
		"CustomField3":         r.FormValue("CustomField3"),
		"CustomField4":         r.FormValue("CustomField4"),
		"MerchantID":           r.FormValue("MerchantID"),
		"MerchantTradeNo":      r.FormValue("MerchantTradeNo"),
		"PaymentDate":          r.FormValue("PaymentDate"),
		"PaymentType":          r.FormValue("PaymentType"),
		"PaymentTypeChargeFee": r.FormValue("PaymentTypeChargeFee"),
		"RtnCode":              r.FormValue("RtnCode"),
		"RtnMsg":               r.FormValue("RtnMsg"),
		"SimulatePaid":         r.FormValue("SimulatePaid"),
		"StoreID":              r.FormValue("StoreID"),
		"TradeAmt":             r.FormValue("TradeAmt"),
		"TradeDate":            r.FormValue("TradeDate"),
		"TradeNo":              r.FormValue("TradeNo"),
		"CheckMacValue":        r.FormValue("CheckMacValue"),
	}

	// 存儲前進行 CheckMacValue 的 Base64 編碼
	if checkMacValue, exists := orderInfo["CheckMacValue"]; exists {
		// 將原始 CheckMacValue 進行 Base64 編碼
		encodedCheckMac := base64.StdEncoding.EncodeToString([]byte(checkMacValue))
		// 保存原始值到一個特殊的鍵中
		orderInfo["CheckMacValue"] = encodedCheckMac
	}

	dataJSON, err := json.Marshal(orderInfo)
	if err != nil {
		log.Printf("JSON Marshal Error: %v", err)
		return
	}

	// 更新 Redis
	key := fmt.Sprintf("order:%s", orderInfo["MerchantTradeNo"])
	if err := c.redisClient.Set(key, dataJSON, 15*time.Minute).Err(); err != nil {
		log.Printf("Redis Set Error: %v", err)
	}

	// 處理回應
	if orderInfo["RtnCode"] == "1" {
		// 交易成功
		http.Redirect(w, r, fmt.Sprintf("/Venue/OrderPending?merchantTradeNo=%s", orderInfo["MerchantTradeNo"]), http.StatusSeeOther)
	} else {
		// 交易失敗
		failOrder := venuepage.OrderPending{
			IsPayFail: true,
		}
		c.RenderTemplate(w, r, "order_pending", failOrder)
	}
}
