package services

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	repoInterfaces "rentjoy/internal/interfaces/repositories"
	serviceInterfaces "rentjoy/internal/interfaces/services"
	"rentjoy/internal/models"
	"rentjoy/internal/repositories"
	"rentjoy/pkg/helper"
	"sort"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EcpayService struct {
	venueInfoRepo repoInterfaces.VenueInformationRepository
	db            *gorm.DB
	merchantID    string
	hashKey       string
	hashIV        string
}

func NewEcpayService(db *gorm.DB) serviceInterfaces.EcpayService {
	return &EcpayService{
		venueInfoRepo: repositories.NewVenueInformationRepository(db),
		db:            db,
		merchantID:    os.Getenv("ECPAY_MERCHANT_ID"), // 測試商店編號
		hashKey:       os.Getenv("ECPAY_HASH_KEY"),    // 測試 HashKey
		hashIV:        os.Getenv("ECPAY_HASH_IV"),     // 測試 HashIV
	}
}

// 包裝 Order 資料成傳給 Ecpay 第三方的資料格式
func (s *EcpayService) PostOrderDetails(order *models.Order, r *http.Request) (map[string]string, error) {
	// 從請求中獲取 baseURL
	baseURL := fmt.Sprintf("https://%s%s", r.Host, r.URL.Port())
	// 生成訂單編號
	orderUUID := s.generateOrderID()

	// 取得場地
	venue, err := s.venueInfoRepo.FindByID(order.VenueID)
	if err != nil {
		return nil, err
	}

	// 組裝 ECPay 需要的參數
	ecpayParams := map[string]string{
		"MerchantTradeNo":   orderUUID,
		"MerchantTradeDate": order.CreatedAt.Format("2006/01/02 15:04:05"),
		"TotalAmount":       order.Amount.StringFixed(0),
		"TradeDesc":         "無",
		"ItemName":          venue.Name,
		"ReturnURL":         fmt.Sprintf("%s/Ecpay/ReceivePaymentResult", baseURL),
		"OrderResultURL":    fmt.Sprintf("%s/Ecpay/ReceivePaymentResult", baseURL),
		"MerchantID":        s.merchantID,
		"PaymentType":       "aio",
		"ChoosePayment":     "Credit",
		"EncryptType":       "1",
	}

	// 計算檢查碼
	ecpayParams["CheckMacValue"] = s.getCheckMacValue(ecpayParams)

	return ecpayParams, nil
}

// 解析 Ecpay 回傳資料
func (s *EcpayService) GetOrderDetails(orderInfo map[string]string) map[string]string {
	ecpayReturnData := make(map[string]string)
	for key, value := range orderInfo {
		if key == "CheckMacValue" {
			continue
		}

		if value == "" {
			ecpayReturnData[key] = ""
		} else {
			ecpayReturnData[key] = value
		}
	}

	// 計算並加入 CheckMacValue
	ecpayReturnData["CheckMacValue"] = s.getCheckMacValue(ecpayReturnData)

	return ecpayReturnData
}

// 計算檢查碼
func (s *EcpayService) getCheckMacValue(params map[string]string) string {
	// 按照鍵值排序
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 組合參數字串
	var values []string
	for _, k := range keys {
		values = append(values, k+"="+params[k])
	}
	checkString := "HashKey=" + s.hashKey + "&" + strings.Join(values, "&") + "&HashIV=" + s.hashIV

	// URL encode
	checkString = url.QueryEscape(checkString)
	checkString = strings.ToLower(checkString)

	// SHA256 加密
	return strings.ToUpper(helper.GetSHA256(checkString))
}

// 產生訂單 uuid
func (s *EcpayService) generateOrderID() string {
	uuid := uuid.New().String()
	// 移除所有非字母和數字的字符
	uuid = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1 // 移除該字符
	}, uuid)

	// 限制長度，確保不超過 ECPay 的限制
	if len(uuid) > 18 {
		uuid = uuid[:18]
	}

	return "RT" + uuid
}
