package manage

// 預訂單管理頁面
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

// 訂單資訊
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

// 場地管理頁面
type VenueManagement struct {
	PublishedVenues  []VenueInfo `json:"publishedVenues"`
	RejectedVenues   []VenueInfo `json:"rejectedVenues"`
	ProcessingVenues []VenueInfo `json:"processingVenues"`
	DelistVenues     []VenueInfo `json:"editingVenues"`
}

type VenueInfo struct {
	VenueId      string `json:"venueId"`
	VenueName    string `json:"venueName"`
	VenueManager string `json:"venueManager"`
	ImgUrl       string `json:"imgUrl"`
}
