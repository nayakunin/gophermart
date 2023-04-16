package storage

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUserExists               = errors.New("user exists")
	ErrUserNotFound             = errors.New("user not found")
	ErrWithdrawOrderNotFound    = errors.New("withdraw order not found")
	ErrWithdrawBalanceNotEnough = errors.New("withdraw balance not enough")
)

func (s *DBStorage) CreateUser(email, password string) error {
	_, err := s.Pool.Exec(context.Background(), `INSERT INTO users (email, password) VALUES ($1, $2)`, email, password)
	if err != nil {
		return err
	}

	return nil
}

func (s *DBStorage) ValidateUser(email, password string) error {
	var userID string
	err := s.Pool.QueryRow(context.Background(), `SELECT id FROM users WHERE email = $1 AND password = $2`, email, password).Scan(&userID)
	if err != nil {
		return ErrUserNotFound
	}

	return nil
}

func (s *DBStorage) SaveOrder(orderID int64, userID string) error {
	_, err := s.Pool.Exec(context.Background(), `INSERT INTO orders (id, user_id) VALUES ($1, $2, $3)`, orderID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *DBStorage) GetOrders(userID string) ([]Order, error) {
	rows, err := s.Pool.Query(context.Background(), `SELECT (id, uploaded_at) FROM orders WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}

	orders := make([]Order, 0)
	for rows.Next() {
		var orderID int64
		var uploadedAt time.Time
		err := rows.Scan(&orderID, &uploadedAt)
		if err != nil {
			return nil, err
		}

		orders = append(orders, Order{
			ID:         orderID,
			UploadedAt: uploadedAt,
		})
	}

	return orders, nil
}

func (s *DBStorage) GetBalance(userID string) (balance float32, withdrawn float32, err error) {
	err = s.Pool.QueryRow(context.Background(), `SELECT (amount, withdrawn) FROM balances WHERE user_id = $1`, userID).Scan(&balance, &withdrawn)

	return balance, withdrawn, err
}

func (s *DBStorage) Withdraw(userID string, order int64, amount float32) error {
	conn, _ := s.Pool.Acquire(context.Background())

	// Check if order exists
	_, err := conn.Exec(context.Background(), `SELECT id FROM orders WHERE id = $1 AND user_id = $2`, order, userID)
	if err != nil {
		return ErrWithdrawOrderNotFound
	}

	// Check if balance is enough
	var balance float32
	err = conn.QueryRow(context.Background(), `SELECT amount FROM balances WHERE user_id = $1`, userID).Scan(&balance)
	if err != nil {
		return err
	}
	if balance < amount {
		return ErrWithdrawBalanceNotEnough
	}

	// Withdraw
	_, err = conn.Exec(context.Background(), `UPDATE balances SET withdrawn = withdrawn + $1, amount = amount - $1 WHERE user_id = $2`, amount, userID)
	if err != nil {
		return err
	}

	// Save transaction
	_, err = conn.Exec(context.Background(), `INSERT INTO transactions (user_id, amount) VALUES ($1, $2)`, userID, amount)
	if err != nil {
		return err
	}

	conn.Release()

	return nil
}

func (s *DBStorage) GetWithdrawals(userID string) ([]Transaction, error) {
	rows, err := s.Pool.Query(context.Background(), `SELECT (order_id, amount, processed_at) FROM transactions WHERE user_id = $1 ORDER BY processed_at ASC`, userID)
	if err != nil {
		return nil, err
	}

	transactions := make([]Transaction, 0)
	for rows.Next() {
		var orderID int64
		var amount float32
		var processedAt time.Time
		err := rows.Scan(&orderID, &amount, &processedAt)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, Transaction{
			OrderID:     orderID,
			Amount:      amount,
			ProcessedAt: processedAt,
		})
	}

	return transactions, nil
}
