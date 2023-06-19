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

	t, err := conn.Begin(context.Background())
	if err != nil {
		return 0, err
	}

	var userID int64
	err = t.QueryRow(context.Background(), `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`, email, password).Scan(&userID)
	if err != nil {
		return 0, err
	}

	_, err = t.Exec(context.Background(), `INSERT INTO balances (user_id, amount, withdrawn) VALUES ($1, 0, 0)`, userID)
	if err != nil {
		return 0, err
	}

	err = t.Commit(context.Background())
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

func (s *DBStorage) SaveOrder(userID, orderID int64) error {
	conn, err := s.Pool.Acquire(context.Background())
	if err != nil {
		return err
	}

	t, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}

	res, err := t.Exec(context.Background(), `INSERT INTO orders (order_id, user_id, status) VALUES ($1, $2, $3) ON CONFLICT (order_id) DO NOTHING`, orderID, userID, api.OrderStatusNEW.ToValue())
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		var ownerID int64
		err := t.QueryRow(context.Background(), `SELECT user_id FROM orders WHERE order_id = $1`, orderID).Scan(&ownerID)
		if err != nil {
			return err
		}

		if ownerID != userID {
			return ErrSaveOrderConflict
		}

		return ErrSaveOrderAlreadyExists
	}

	err = t.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (s *DBStorage) GetOrders(userID int64) ([]Order, error) {
	rows, err := s.Pool.Query(context.Background(), `SELECT (order_id, status, uploaded_at, accrual) FROM orders WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}

	orders := make([]Order, 0)
	for rows.Next() {
		type orderRow struct {
			OrderID    int64
			Status     string
			UploadedAt time.Time
			Accrual    *float32
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
			Accrual:    row.Accrual,
			Status:     status,
			UploadedAt: row.UploadedAt,
		})
	}

	return orders, nil
}

func (s *DBStorage) GetBalance(userID int64) (Balance, error) {
	var balance Balance
	err := s.Pool.QueryRow(context.Background(), `SELECT (amount, withdrawn) FROM balances WHERE user_id = $1`, userID).Scan(&balance)
	if err != nil {
		return balance, err
	}

	return balance, err
}

func (s *DBStorage) Withdraw(userID, orderID int64, amount float32) error {
	conn, err := s.Pool.Acquire(context.Background())
	if err != nil {
		return err
	}

	defer conn.Release()

	t, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}

	// Check if orderID exists
	_, err = t.Exec(context.Background(), `SELECT id FROM orders WHERE order_id = $1 AND user_id = $2`, orderID, userID)
	if err != nil {
		return ErrWithdrawOrderNotFound
	}

	// Check if balance is enough
	var balance float32
	err = t.QueryRow(context.Background(), `SELECT amount FROM balances WHERE user_id = $1`, userID).Scan(&balance)
	if err != nil {
		return err
	}
	if balance < amount {
		return ErrWithdrawBalanceNotEnough
	}

	// Withdraw
	_, err = t.Exec(context.Background(), `UPDATE balances SET withdrawn = withdrawn + $1, amount = amount - $1 WHERE user_id = $2`, amount, userID)
	if err != nil {
		return err
	}

	// Save transaction
	_, err = t.Exec(context.Background(), `INSERT INTO transactions (user_id, order_id, amount) VALUES ($1, $2, $3)`, userID, orderID, amount)
	if err != nil {
		return err
	}

	err = t.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (s *DBStorage) GetWithdrawals(userID int64) ([]Transaction, error) {
	rows, err := s.Pool.Query(context.Background(), `SELECT (order_id, amount, processed_at) FROM transactions WHERE user_id = $1 ORDER BY processed_at ASC`, userID)
	if err != nil {
		return nil, err
	}

	transactions := make([]Transaction, 0)
	for rows.Next() {
		type Row struct {
			OrderID     int64
			Amount      float32
			ProcessedAt time.Time
		}
		var row Row
		err := rows.Scan(&row)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, Transaction{
			OrderID:     row.OrderID,
			Amount:      row.Amount,
			ProcessedAt: row.ProcessedAt,
		})
	}

	return transactions, nil
}

func (s *DBStorage) UpdateOrderStatus(orderID int64, status api.OrderStatus) error {
	_, err := s.Pool.Exec(context.Background(), `UPDATE orders SET status = $1 WHERE order_id = $2`, status.ToValue(), orderID)
	if err != nil {
		return err
	}

	return nil
}

func (s *DBStorage) ProcessOrder(userID int64, orderID int64, accrual float32) error {
	conn, err := s.Pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	t, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}

	_, err = t.Exec(context.Background(), `UPDATE orders SET accrual = $1, status = $2 WHERE order_id = $3;`, accrual, api.OrderStatusPROCESSED.ToValue(), orderID)
	if err != nil {
		return err
	}

	_, err = t.Exec(context.Background(), `UPDATE balances SET amount = amount + $1 WHERE user_id = $2`, accrual, userID)
	if err != nil {
		return err
	}

	err = t.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}
