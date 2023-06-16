package server

import (
	"github.com/nayakunin/gophermart/internal/storage"
)

type Storage interface {
	CreateUser(email string, password string) (int64, error)
	GetUserID(email string, password string) (int64, error)
	SaveOrder(userID int64, orderID int64) error
	GetOrders(userID int64) ([]storage.Order, error)
	GetBalance(userID int64) (float32, float32, error)
	Withdraw(userID int64, orderID int64, amount float32) error
	GetWithdrawals(userID int64) ([]storage.Transaction, error)
}

type Worker interface {
	AddOrder(userID int64, orderID int64)
	Start()
}
