package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strconv"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// 字串雜湊
func GetSHA256(p string) string {
	hasher := sha256.New()

	hasher.Write([]byte(p))

	hashedBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashedBytes)
}

// 產生隨機密碼
func GenerateRandomPassword() string {
	return uuid.New().String()
}

// 輸入密碼比對
func CheckPasswordHash(inputPassword, modelPassword string) bool {
	if GetSHA256(inputPassword) != modelPassword {
		return false
	} else {
		return true
	}

}

// 正則表達式驗證手機號碼格式
func ValidatePhoneNumber(phone string) bool {
	phoneRegex := regexp.MustCompile(`^09\d{8}$`)
	return phoneRegex.MatchString(phone)
}

// string 轉 uint
func StrToUint(str string) (uint, error) {
	if str == "" {
		return 0, nil
	}

	num, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(num), nil
}

// string 轉 int
func StringToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	return num
}

// string slice 轉 int slice
func StringSliceToIntSlice(strSlice []string) []int {
	intSlice := make([]int, len(strSlice))
	for i, str := range strSlice {
		intSlice[i] = StringToInt(str)
	}
	return intSlice
}

// decimal 轉 int
func DecimalToIntRounded(d decimal.Decimal) int {
	return int(d.Round(0).IntPart())
}

// 根據查詢條件取得對應的小時範圍
func GetTimeSlotCondition(timeSlot string) (startHour, endHour string) {
	switch timeSlot {
	case "上午":
		return "00:00:00", "11:59:59"
	case "下午":
		return "12:00:00", "17:59:59"
	case "晚上":
		return "18:00:00", "23:59:59"
	default:
		return "00:00:00", "23:59:59"
	}
}

// 根據查詢條件取得對應的星期範圍
func GetDayTypeCondition(dayType string) []string {
	switch dayType {
	case "平日":
		return []string{"1", "2", "3", "4", "5"} // Monday to Friday
	case "假日":
		return []string{"0", "6"} // Sunday and Saturday
	default:
		return []string{"0", "1", "2", "3", "4", "5", "6"}
	}
}
