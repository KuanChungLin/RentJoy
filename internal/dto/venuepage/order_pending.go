package venuepage

type OrderPending struct {
	VenueId             string `json:"venueId"`
	OrderId             string `json:"orderId"`
	OrderNo             string `json:"orderNo"`
	Email               string `json:"email"`
	Phone               string `json:"phone"`
	IsPayFail           bool   `json:"isPayFail"`
	IsCreateFail        bool   `json:"isCreateFail"`
	IsCheckMacValueFail bool   `json:"isCheckMacValueFail"`
}
