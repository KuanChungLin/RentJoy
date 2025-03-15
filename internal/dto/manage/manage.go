package manage

type ReservedManagement struct {
	OrderCount    int         `json:"orderCount"`
	AcceptCount   int         `json:"acceptCount"`
	BookingAmount int         `json:"bookingAmount"`
	RejectCount   int         `json:"rejectCount"`
	PendingCount  int         `json:"pendingCount"`
	Orders        []OrderInfo `json:"orders"`
	CurrentPage   int         `json:"currentPage"`
	TotalPages    int         `json:"totalpages"`
}

type OrderInfo struct {
	OrderId     string `json:"orderId"`
	OrderDesc   string `json:"orderDesc"`
	VenueName   string `json:"venueName"`
	Booker      string `json:"booker"`
	BookingTime string `json:"bookingTime"`
	Phone       string `json:"phone"`
	Amount      int    `json:"amount"`
	Status      string `json:"status"`
	OrderTime   string `json:"orderTime"`
}
