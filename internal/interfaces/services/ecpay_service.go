package interfaces

import (
	"net/http"
	"rentjoy/internal/models"
)

type EcpayService interface {
	PostOrderDetails(order *models.Order, r *http.Request) (map[string]string, error)
	GetOrderDetails(orderInfo map[string]string) map[string]string
}
