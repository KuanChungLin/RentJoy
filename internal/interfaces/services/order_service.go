package interfaces

import (
	"net/http"
	"rentjoy/internal/dto/order"
)

type OrderService interface {
	SaveOrder(orderForm order.OrderForm, userID uint, r *http.Request) (map[string]string, string, error)
}
