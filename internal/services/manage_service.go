package services

import (
	"log"
	"rentjoy/internal/dto/manage"
	"rentjoy/internal/dto/order"
	repoInterfaces "rentjoy/internal/interfaces/repositories"
	serviceInterfaces "rentjoy/internal/interfaces/services"
	"rentjoy/internal/repositories"
	"rentjoy/pkg/helper"
	"strconv"

	"gorm.io/gorm"
)

type ManageService struct {
	orderRepo repoInterfaces.OrderRepository
	venueRepo repoInterfaces.VenueInformationRepository
}

func NewManageService(db *gorm.DB) serviceInterfaces.ManageService {
	return &ManageService{
		orderRepo: repositories.NewOrderRepository(db),
		venueRepo: repositories.NewVenueInformationRepository(db),
	}
}

// 取得預訂單管理頁面資訊
func (s *ManageService) GetReservedManagement(userId uint) (*manage.ReservedManagement, error) {
	orderInfos, err := s.orderRepo.FindManageOrderByUserId(userId)
	if err != nil {
		return &manage.ReservedManagement{}, err
	}

	// 宣告 ReservedManagement 所需欄位
	var orders []manage.OrderInfo
	var acceptCount, rejectCount, pendingCount, bookingAmount int

	for _, oi := range orderInfos {
		var orderStatus string
		detail := oi.Details[0]
		orderStatus = order.OrderStatus(oi.Status).String()

		switch oi.Status {
		case 1:
			pendingCount++
		case 2:
			acceptCount++
			bookingAmount += helper.DecimalToIntRounded(oi.Amount)
		case 3:
			rejectCount++
		}

		order := manage.OrderInfo{
			OrderId:     strconv.Itoa(int(oi.ID)),
			OrderDesc:   oi.Message,
			VenueName:   oi.Venue.Name,
			Booker:      oi.LastName + oi.FirstName,
			BookingTime: detail.StartTime.Format("2006-01-02 12:30"),
			Phone:       oi.Phone,
			Amount:      helper.DecimalToIntRounded(oi.Amount),
			Status:      orderStatus,
			OrderTime:   oi.CreatedAt.Format("2006-01-02 12:30"),
		}

		orders = append(orders, order)
	}

	reserved := manage.ReservedManagement{
		OrderCount:    len(orderInfos),
		AcceptCount:   acceptCount,
		BookingAmount: bookingAmount,
		RejectCount:   rejectCount,
		PendingCount:  pendingCount,
		Orders:        orders,
	}

	return &reserved, nil
}

func (s *ManageService) GetVenueManagement(ownerId uint) (*manage.VenueManagement, error) {
	venues, err := s.venueRepo.FindByOwnerId(ownerId)
	if err != nil {
		log.Printf("Get Venue Management By OwnerId Error:%s", err)
		return &manage.VenueManagement{}, err
	}

	var publishedVenues []manage.VenueInfo
	var rejectedVenues []manage.VenueInfo
	var processingVenues []manage.VenueInfo
	var delistVenues []manage.VenueInfo

	for _, venue := range venues {
		venueinfo := manage.VenueInfo{
			VenueId:      strconv.Itoa(int(venue.ID)),
			VenueName:    venue.Name,
			VenueManager: venue.Owner.LastName + venue.Owner.FirstName,
			ImgUrl:       venue.Imgs[0].VenueImgPath,
		}

		switch venue.Status {
		case 1:
			publishedVenues = append(publishedVenues, venueinfo)
		case 2:
			processingVenues = append(processingVenues, venueinfo)
		case 3:
			delistVenues = append(delistVenues, venueinfo)
		case 4:
			rejectedVenues = append(rejectedVenues, venueinfo)

		}
	}

	vm := manage.VenueManagement{
		PublishedVenues:  publishedVenues,
		RejectedVenues:   rejectedVenues,
		ProcessingVenues: processingVenues,
		DelistVenues:     delistVenues,
	}

	return &vm, nil
}

// 接受場地預訂作業
func (s *ManageService) ReservedAccept(orderId uint) bool {
	order, err := s.orderRepo.FindByID(orderId)
	if err != nil {
		log.Printf("Get Order By OrderId Error:%s", err)
		return false
	}

	order.Status = 2

	err = s.orderRepo.Update(*order)
	if err != nil {
		log.Printf("Update Order Error:%s", err)
		return false
	}
	return true
}

// 拒絕場地預訂作業
func (s *ManageService) ReservedReject(orderId uint) bool {
	order, err := s.orderRepo.FindByID(orderId)
	if err != nil {
		log.Printf("Get Order By OrderId Error:%s", err)
		return false
	}

	order.Status = 3

	err = s.orderRepo.Update(*order)
	if err != nil {
		log.Printf("Update Order Error:%s", err)
		return false
	}
	return true
}
