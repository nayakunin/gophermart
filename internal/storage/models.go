package storage

import "time"

type User struct {
	ID       int64
	Email    string
	Password string
}

type Order struct {
	ID         int64
	UploadedAt time.Time
}

type Balance struct {
	ID        int64
	UserID    int64
	Amount    int64
	Withdrawn int64
}

type Transaction struct {
	ID          int64
	UserID      int64
	OrderID     int64
	Amount      int64
	ProcessedAt time.Time
}
