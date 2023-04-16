package storage

import (
	"time"

	api "github.com/nayakunin/gophermart/internal/generated"
)

type User struct {
	ID       int64
	Email    string
	Password string
}

type Order struct {
	ID         int64
	Status     api.OrderStatus
	UploadedAt time.Time
}

type Balance struct {
	ID        int64
	UserID    int64
	Amount    float32
	Withdrawn float32
}

type Transaction struct {
	ID          int64
	UserID      int64
	OrderID     int64
	Amount      float32
	ProcessedAt time.Time
}
