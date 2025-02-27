package order

type OrderForm struct {
	Activity        string `json:"activity"`
	UserCount       string `json:"userCount"`
	Message         string `json:"message"`
	LastName        string `json:"lastName"`
	FirstName       string `json:"firstName"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	VenueID         string `json:"venueId"`
	ReservedDetails string `json:"reservedDetails"`
}
