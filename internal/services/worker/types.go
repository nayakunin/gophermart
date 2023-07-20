package worker

import (
	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/nayakunin/gophermart/internal/services/accrual"
)

type AccrualService interface {
	GetAccrual(orderID int64) (*accrual.Accrual, error)
}

type Storage interface {
	UpdateOrderStatus(orderID int64, status api.OrderStatus) error
	ProcessOrder(userID int64, orderID int64, accrual float32) error
}
