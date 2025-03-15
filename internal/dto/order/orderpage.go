package order

import (
	"rentjoy/internal/dto/venuepage"
)

type OrderPageInfo struct {
	Orders        []OrderVM               `json:"orders"`
	Recommend     []venuepage.Recommended `json:"recommend"`
	OrderCount    int                     `json:"orderCount"`
	TotalPages    int                     `json:"totalPages"`
	CurrentPage   int                     `json:"currentPage"`
	CurrentAction string                  `json:"currentAction"`
}

type OrderVM struct {
	Title          string              `json:"title"`
	Address        string              `json:"address"`
	OrderId        string              `json:"orderId"`
	OrderTime      string              `json:"orderTime"`
	CancelTime     string              `json:"cancelTime"`
	Status         OrderStatus         `json:"status"`
	OrderPrice     int                 `json:"orderPrice"`
	ContactPerson  string              `json:"contactPerson"`
	Email          string              `json:"email"`
	ImgUrl         string              `json:"imgUrl"`
	ProductUrl     string              `json:"productUrl"`
	Evaluate       OrderEvaluate       `json:"evaluate"`
	ScheduledTimes []OrderScheduleTime `json:"scheduleTimes"`
}

type OrderEvaluate struct {
	Rating       int    `json:"rating"`
	Content      string `json:"content"`
	EvaluateTime string `json:"evaluateTime"`
}

type OrderScheduleTime struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type OrderStatus int

const (
	// Unpaid 送出訂單未付款
	Unpaid OrderStatus = iota
	// Reserved 處理中(已付款店家未確認)
	Reserved
	// Processing 已預訂(已付款店家已確認)
	Processing
	// Cancel 取消預訂
	Cancel
	// Finished 已結束
	Finished
	// Fail 訂單失效(超過時間or付款失敗)
	Fail
)

func (s OrderStatus) String() string {
	switch s {
	case Unpaid:
		return "未付款"
	case Reserved:
		return "處理中"
	case Processing:
		return "已預訂"
	case Cancel:
		return "取消預訂"
	case Finished:
		return "已結束"
	case Fail:
		return "訂單失效"
	default:
		return "未知狀態"
	}
}
