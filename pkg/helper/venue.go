package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"rentjoy/internal/dto/venuepage"
	"rentjoy/internal/models"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// 依照場地圖片的 Sort 進行排序
func GetSortedImgs(imgs []models.VenueImg) []string {
	sort.Slice(imgs, func(i, j int) bool {
		return imgs[i].Sort < imgs[j].Sort
	})

	paths := make([]string, len(imgs))

	for i, img := range imgs {
		paths[i] = img.VenueImgPath
	}

	return paths
}

// 將 model 訂單及評價資料轉換成 dto
func GetVenueComments(orders []models.Order) []venuepage.Comment {
	var comments []venuepage.Comment
	for _, order := range orders {
		if order.VenueEvaluate != nil && order.Status == 4 {
			comments = append(comments, venuepage.Comment{
				UserName:     order.FirstName + order.LastName,
				CommentYear:  strconv.Itoa(order.VenueEvaluate.CreatedAt.Year()),
				CommentMonth: strconv.Itoa(int(order.VenueEvaluate.CreatedAt.Month())),
				CommentDay:   strconv.Itoa(order.VenueEvaluate.CreatedAt.Day()),
				CommentTxt:   order.VenueEvaluate.EvaluateComment,
			})
		}
	}
	return comments
}

// 將 model 設備資料轉換成 dto
func ConvertToVenueDevice(devices []models.VenueDevice) []venuepage.VenueDevice {
	if len(devices) == 0 {
		return []venuepage.VenueDevice{} // 返回空陣列
	}

	result := make([]venuepage.VenueDevice, len(devices))
	for i, device := range devices {
		result[i] = venuepage.VenueDevice{
			DeviceName:     device.DeviceItem.DeviceName,
			DeviceQuantity: device.Count,
			DeviceRemark:   device.DeviceDescription,
		}
	}
	return result
}

// 取得未包含的設備清單
func GetNotIncludedDevices(allDevices []string, venueDevices []models.VenueDevice) []string {
	existingDevices := make(map[string]bool)
	for _, device := range venueDevices {
		existingDevices[device.DeviceItem.DeviceName] = true
	}

	var notIncluded []string
	for _, device := range allDevices {
		if !existingDevices[device] {
			notIncluded = append(notIncluded, device)
			if len(notIncluded) >= 18 {
				break
			}
		}
	}
	return notIncluded
}

// 將規則字串分割成陣列
func SplitRules(rules string) []string {
	if rules == "" {
		return []string{}
	}
	return strings.Split(rules, "\n")
}

// 取得管理者資訊
func GetOwnerInfo(management *models.Management) venuepage.OwnerInfo {
	if management == nil {
		return venuepage.OwnerInfo{
			ImgUrl:    "/static/images/empty-photo.svg",
			Name:      "未設定管理者名稱",
			JoinYear:  "",
			JoinMonth: "",
			JoinDay:   "",
		}
	}

	return venuepage.OwnerInfo{
		ImgUrl:    management.AvatarImgLinkPath,
		Name:      management.ManagementName,
		JoinYear:  strconv.Itoa(management.CreatedAt.Year()),
		JoinMonth: strconv.Itoa(int(management.CreatedAt.Month())),
		JoinDay:   strconv.Itoa(management.CreatedAt.Day()),
	}
}

// 取得每小時金額及每時段金額
func GetPriceRange(rateTypeId uint, billingRates []models.BillingRate) string {
	var minRate decimal.Decimal
	var maxRate decimal.Decimal
	var hasRate bool

	for _, rate := range billingRates {
		if rate.RateTypeID == rateTypeId {
			if !hasRate {
				minRate = rate.Rate
				maxRate = rate.Rate
				hasRate = true
			}

			if rate.Rate.LessThan(minRate) {
				minRate = rate.Rate
			}
			if rate.Rate.GreaterThan(maxRate) {
				maxRate = rate.Rate
			}
		}
	}

	if !hasRate {
		return "尚未設定價格"
	}

	if minRate.Equal(maxRate) {
		return fmt.Sprintf("$ %s", minRate.StringFixed(0))
	}

	return fmt.Sprintf("$ %s - $ %s", minRate.StringFixed(0), maxRate.StringFixed(0))
}

// 取得預約日期
func GetReserveDates(orders []models.Order) []string {
	dateMap := make(map[time.Time]bool)
	var dates []string

	for _, order := range orders {
		if order.Status >= 0 && order.Status <= 2 {
			for _, detail := range order.Details {
				date := detail.StartTime.Truncate(24 * time.Hour)
				if !dateMap[date] {
					dateMap[date] = true
					dates = append(dates, date.String())
				}
			}
		}
	}

	return dates
}

// 取得唯一的星期幾
func GetUniqueDayOfWeek(rates []models.BillingRate) []int {
	dayMap := make(map[time.Weekday]bool)
	var days []int

	for _, rate := range rates {
		if !dayMap[rate.DayOfWeek] {
			dayMap[rate.DayOfWeek] = true
			days = append(days, int(rate.DayOfWeek))
		}
	}

	return days
}

// 取得最小租用時數
func GetMinRentHours(rates []models.BillingRate) float32 {
	for _, rate := range rates {
		if rate.RateTypeID == 1 {
			return rate.MinRentHours
		}
	}
	return 0
}

// 檢查日期是否有效 (從後天開始預約場地)
func IsDateValid(date time.Time) bool {
	now := time.Now()
	minValidDate := now.Add(24 * time.Hour)
	return date.After(minValidDate)
}

// 標準化時間，只保留時分秒
func NormalizeTime(t time.Time) time.Time {
	return time.Date(0, 1, 1, t.Hour(), t.Minute(), 0, 0, t.Location())
}

// 檢查時間是否衝突
func IsTimeConflict(start, end time.Time, orders []models.Order) bool {
	normalizedStart := NormalizeTime(start)
	normalizedEnd := NormalizeTime(end)

	for _, order := range orders {
		for _, detail := range order.Details {
			orderStart := NormalizeTime(detail.StartTime)
			orderEnd := NormalizeTime(detail.EndTime)

			if normalizedStart.Equal(orderStart) ||
				(normalizedStart.After(orderStart) && normalizedStart.Before(orderEnd)) ||
				(normalizedEnd.After(orderStart) && normalizedEnd.Before(orderEnd)) {
				return true
			}
		}
	}
	return false
}

// 將星期從數字轉換成中文
func GetDayOfWeekInChinese(date time.Time) string {
	switch date.Weekday() {
	case time.Sunday:
		return "星期日"
	case time.Monday:
		return "星期一"
	case time.Tuesday:
		return "星期二"
	case time.Wednesday:
		return "星期三"
	case time.Thursday:
		return "星期四"
	case time.Friday:
		return "星期五"
	case time.Saturday:
		return "星期六"
	default:
		return ""
	}
}

// FormatAddress 格式化地址
func FormatAddress(city, district, address string) string {
	// 移除多餘的空格
	city = strings.TrimSpace(city)
	district = strings.TrimSpace(district)
	address = strings.TrimSpace(address)
	log.Println("city", city)
	log.Println("district", district)
	log.Println("address", address)

	// 按照正确的顺序组合地址
	fullAddress := fmt.Sprintf("%s%s%s", city, district, address)
	log.Println("fullAddress", fullAddress)

	// 移除可能存在的重複空格
	fullAddress = strings.Join(strings.Fields(fullAddress), " ")
	log.Println("fullAddress after trim", fullAddress)
	return fullAddress
}

// GetCoordinates 取得地址的經緯度
func GetCoordinates(city, district, address string) (string, string) {
	// 格式化地址
	fullAddress := FormatAddress(city, district, address)

	// 從環境變量獲取 Google Maps API Key
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		log.Println("Error: GOOGLE_MAPS_API_KEY is not set")
		return "", ""
	}

	// 構建請求 URL
	baseURL := "https://maps.googleapis.com/maps/api/geocode/json"
	params := url.Values{}
	params.Add("address", fullAddress)
	params.Add("key", apiKey)
	params.Add("region", "tw")      // 指定台灣地區
	params.Add("language", "zh-TW") // 指定語言為繁體中文

	requestURL := baseURL + "?" + params.Encode()

	// 發送請求
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Error getting coordinates: %v\n", err)
		return "", ""
	}
	defer resp.Body.Close()

	// 解析響應
	var result struct {
		Status  string
		Results []struct {
			Geometry struct {
				Location struct {
					Lat float64
					Lng float64
				}
			}
		}
		ErrorMessage string `json:"error_message,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decoding response: %v\n", err)
		return "", ""
	}

	// 檢查響應狀態
	if result.Status != "OK" {
		log.Printf("Geocoding failed with status: %s, error message: %s\n", result.Status, result.ErrorMessage)
		return "", ""
	}

	if len(result.Results) == 0 {
		log.Println("No results found for the given address")
		return "", ""
	}

	// 取得經緯度
	location := result.Results[0].Geometry.Location
	return fmt.Sprintf("%f", location.Lat), fmt.Sprintf("%f", location.Lng)
}
