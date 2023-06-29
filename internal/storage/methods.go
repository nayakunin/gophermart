package storage

import (
	"context"
	"time"

	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/pkg/errors"
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
		return 0, errors.Wrap(err, "acquire connection")
	}
	defer conn.Release()

	t, err := conn.Begin(context.Background())
	if err != nil {
		return 0, errors.Wrap(err, "begin transaction")
	}

	var userID int64
	if err = t.QueryRow(context.Background(), `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`, email, password).Scan(&userID); err != nil {
		return 0, errors.Wrap(err, "insert user")
	}

	if _, err = t.Exec(context.Background(), `INSERT INTO balances (user_id, amount, withdrawn) VALUES ($1, 0, 0)`, userID); err != nil {
		return 0, errors.Wrap(err, "insert balance")
	}

	if err = t.Commit(context.Background()); err != nil {
		return 0, errors.Wrap(err, "commit transaction")
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
		return errors.Wrap(err, "acquire connection")
	}

	t, err := conn.Begin(context.Background())
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}

	res, err := t.Exec(context.Background(), `INSERT INTO orders (order_id, user_id, status) VALUES ($1, $2, $3) ON CONFLICT (order_id) DO NOTHING`, orderID, userID, api.OrderStatusNEW.ToValue())
	if err != nil {
		return errors.Wrap(err, "insert order")
	}

	if res.RowsAffected() == 0 {
		var ownerID int64
		if err := t.QueryRow(context.Background(), `SELECT user_id FROM orders WHERE order_id = $1`, orderID).Scan(&ownerID); err != nil {
			return errors.Wrap(err, "select order")
		}

		if ownerID != userID {
			return ErrSaveOrderConflict
		}

		return ErrSaveOrderAlreadyExists
	}

	if err = t.Commit(context.Background()); err != nil {
		return errors.Wrap(err, "commit transaction")
	}

	return nil
}

func (s *DBStorage) GetOrders(userID int64) ([]Order, error) {
	rows, err := s.Pool.Query(context.Background(), `SELECT (order_id, status, uploaded_at, accrual) FROM orders WHERE user_id = $1`, userID)
	if err != nil {
		return nil, errors.Wrap(err, "select orders")
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
		if err := rows.Scan(&row); err != nil {
			return nil, errors.Wrap(err, "scan order")
		}

		status := api.OrderStatus{}
		if err = status.FromValue(row.Status); err != nil {
			return nil, errors.Wrap(err, "parse order status")
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
	if err := s.Pool.QueryRow(context.Background(), `SELECT (amount, withdrawn) FROM balances WHERE user_id = $1`, userID).Scan(&balance); err != nil {
		return balance, errors.Wrap(err, "select balance")
	}

	return balance, nil
}

func (s *DBStorage) Withdraw(userID, orderID int64, amount float32) error {
	conn, err := s.Pool.Acquire(context.Background())
	if err != nil {
		return errors.Wrap(err, "acquire connection")
	}

	defer conn.Release()

	t, err := conn.Begin(context.Background())
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}

	// Check if orderID exists
	if _, err = t.Exec(context.Background(), `SELECT id FROM orders WHERE order_id = $1 AND user_id = $2`, orderID, userID); err != nil {
		return ErrWithdrawOrderNotFound
	}

	// Check if balance is enough
	var balance float32
	if err = t.QueryRow(context.Background(), `SELECT amount FROM balances WHERE user_id = $1`, userID).Scan(&balance); err != nil {
		return errors.Wrap(err, "select balance")
	}
	if balance < amount {
		return ErrWithdrawBalanceNotEnough
	}

	// Withdraw
	if _, err = t.Exec(context.Background(), `UPDATE balances SET withdrawn = withdrawn + $1, amount = amount - $1 WHERE user_id = $2`, amount, userID); err != nil {
		return errors.Wrap(err, "update balance")
	}

	// Save transaction
	if _, err = t.Exec(context.Background(), `INSERT INTO transactions (user_id, order_id, amount) VALUES ($1, $2, $3)`, userID, orderID, amount); err != nil {
		return errors.Wrap(err, "insert transaction")
	}

	if err = t.Commit(context.Background()); err != nil {
		return errors.Wrap(err, "commit transaction")
	}

	return nil
}

func (s *DBStorage) GetWithdrawals(userID int64) ([]Transaction, error) {
	rows, err := s.Pool.Query(context.Background(), `SELECT (order_id, amount, processed_at) FROM transactions WHERE user_id = $1 ORDER BY processed_at ASC`, userID)
	if err != nil {
		return nil, errors.Wrap(err, "select transactions")
	}

	transactions := make([]Transaction, 0)
	for rows.Next() {
		type Row struct {
			OrderID     int64
			Amount      float32
			ProcessedAt time.Time
		}
		var row Row
		if err := rows.Scan(&row); err != nil {
			return nil, errors.Wrap(err, "scan transaction")
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
	if _, err := s.Pool.Exec(context.Background(), `UPDATE orders SET status = $1 WHERE order_id = $2`, status.ToValue(), orderID); err != nil {
		return errors.Wrap(err, "update order status")
	}

	return nil
}

func (s *DBStorage) ProcessOrder(userID int64, orderID int64, accrual float32) error {
	conn, err := s.Pool.Acquire(context.Background())
	if err != nil {
		return errors.Wrap(err, "acquire connection")
	}
	defer conn.Release()

	t, err := conn.Begin(context.Background())
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}

	if _, err = t.Exec(context.Background(), `UPDATE orders SET accrual = $1, status = $2 WHERE order_id = $3;`, accrual, api.OrderStatusPROCESSED.ToValue(), orderID); err != nil {
		return errors.Wrap(err, "update order")
	}

	if _, err = t.Exec(context.Background(), `UPDATE balances SET amount = amount + $1 WHERE user_id = $2`, accrual, userID); err != nil {
		return errors.Wrap(err, "update balance")
	}

	if err = t.Commit(context.Background()); err != nil {
		return errors.Wrap(err, "commit transaction")
	}

	return nil
}
