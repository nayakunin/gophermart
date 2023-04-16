package storage

import (
	"context"
	"errors"
	"time"

	api "github.com/nayakunin/gophermart/internal/generated"
)

var (
	ErrUserExists               = errors.New("user exists")
	ErrUserNotFound             = errors.New("user not found")
	ErrWithdrawOrderNotFound    = errors.New("withdraw order not found")
	ErrWithdrawBalanceNotEnough = errors.New("withdraw balance not enough")
	ErrSaveOrderAlreadyExists   = errors.New("save order already exists")
	ErrSaveOrderConflict        = errors.New("save order conflict")
)

func (s *DBStorage) CreateUser(email, password string) (int64, error) {
	conn, err := s.Pool.Acquire(context.Background())
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	var userID int64
	err = conn.QueryRow(context.Background(), `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`, email, password).Scan(&userID)
	if err != nil {
		return 0, err
	}

	_, err = conn.Exec(context.Background(), `INSERT INTO balances (user_id, amount, withdrawn) VALUES ($1, 0, 0)`, userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *DBStorage) GetUserID(email, password string) (int64, error) {
	var userID int64
	err := s.Pool.QueryRow(context.Background(), `SELECT id FROM users WHERE email = $1 AND password = $2`, email, password).Scan(&userID)
	if err != nil {
		return 0, ErrUserNotFound
	}

	return userID, nil
}

func (s *DBStorage) SaveOrder(userID, orderID int64, status string) error {
	res, err := s.Pool.Exec(context.Background(), `INSERT INTO orders (order_id, user_id, status) VALUES ($1, $2, $3) ON CONFLICT (order_id) DO NOTHING`, orderID, userID, status)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		var ownerID int64
		err := s.Pool.QueryRow(context.Background(), `SELECT user_id FROM orders WHERE order_id = $1`, orderID).Scan(&ownerID)
		if err != nil {
			return err
		}

		if ownerID != userID {
			return ErrSaveOrderConflict
		}

		return ErrSaveOrderAlreadyExists
	}

	return nil
}

func (s *DBStorage) GetOrders(userID int64) ([]Order, error) {
	rows, err := s.Pool.Query(context.Background(), `SELECT (order_id, status, uploaded_at) FROM orders WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}

	orders := make([]Order, 0)
	for rows.Next() {
		type orderRow struct {
			OrderID    int64
			Status     string
			UploadedAt time.Time
		}
		var row orderRow
		err := rows.Scan(&row)
		if err != nil {
			return nil, err
		}

		status := api.OrderStatus{}
		err = status.FromValue(row.Status)
		if err != nil {
			return nil, err
		}

		orders = append(orders, Order{
			ID:         row.OrderID,
			Status:     status,
			UploadedAt: row.UploadedAt,
		})
	}

	return orders, nil
}

func (s *DBStorage) GetBalance(userID int64) (float32, float32, error) {
	type balanceRow struct {
		Amount   float32
		Withdraw float32
	}
	var row balanceRow
	err := s.Pool.QueryRow(context.Background(), `SELECT (amount, withdrawn) FROM balances WHERE user_id = $1`, userID).Scan(&row)
	if err != nil {
		return 0, 0, err
	}

	return row.Amount, row.Withdraw, err
}

func (s *DBStorage) Withdraw(userID, order int64, amount float32) error {
	conn, _ := s.Pool.Acquire(context.Background())

	// Check if order exists
	_, err := conn.Exec(context.Background(), `SELECT id FROM orders WHERE order_id = $1 AND user_id = $2`, order, userID)
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

func (s *DBStorage) GetWithdrawals(userID int64) ([]Transaction, error) {
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
