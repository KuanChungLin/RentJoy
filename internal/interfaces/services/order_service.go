package interfaces

import (
	"net/http"
	"rentjoy/internal/dto/order"
)

type OrderService interface {
	SaveOrder(orderForm order.OrderForm, userID uint, r *http.Request) (map[string]string, string, error)
	GetOrderPage(userId uint, status order.OrderStatus, pageIndex int, pageSize int) (*order.OrderPageInfo, error)
	CancelReservation(orderId uint) error
}
