package venuepage

import "rentjoy/internal/dto/homepage"

type ReservedPage struct {
	VenueId          string              `json:"venueId"`
	VenueImgUrl      string              `json:"venueImgUrl"`
	Name             string              `json:"name"`
	Address          string              `json:"address"`
	Date             string              `json:"date"`
	OrderDetails     []Orderdetail       `json:"orderDetails"`
	Amount           string              `json:"amount"`
	ReservedActivity []homepage.Activity `json:"reservedActivity"`
}

type Orderdetail struct {
	TimeDetail string `json:"timeDetail"`
	Price      string `json:"price"`
}
