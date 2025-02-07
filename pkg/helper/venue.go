package helper

import (
	"fmt"
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
