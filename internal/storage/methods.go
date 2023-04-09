package storage

import (
	"context"
	"time"
)

func (s *DBStorage) CreateUser(email, password string) error {
	_, err := s.Pool.Exec(context.Background(), `INSERT INTO users (email, password) VALUES ($1, $2)`, email, password)
	if err != nil {
		return err
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

func (s *DBStorage) GetBalance(userID string) (balance int64, withdrawn int64, err error) {
	err = s.Pool.QueryRow(context.Background(), `SELECT (amount, withdrawn) FROM balances WHERE user_id = $1`, userID).Scan(&balance, &withdrawn)

	return balance, withdrawn, err
}

func (s *DBStorage) Withdraw(userID string, amount int64) error {
	_, err := s.Pool.Exec(context.Background(), `UPDATE balances SET withdrawn = withdrawn - $1, amount = amount + $1 WHERE user_id = $2`, amount, userID)
	if err != nil {
		return err
	}
	_, err = s.Pool.Exec(context.Background(), `INSERT INTO transactions (user_id, amount) VALUES ($1, $2)`, userID, amount)
	if err != nil {
		return err
	}

	return nil
}

func (s *DBStorage) GetWithdrawals(userID string) ([]Transaction, error) {
	rows, err := s.Pool.Query(context.Background(), `SELECT (order_id, amount, processed_at) FROM transactions WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}

	transactions := make([]Transaction, 0)
	for rows.Next() {
		var orderID int64
		var amount int64
		var createdAt time.Time
		err := rows.Scan(&orderID, &amount, &createdAt)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, Transaction{
			OrderID:     orderID,
			Amount:      amount,
			ProcessedAt: createdAt,
		})
	}

	return transactions, nil
}
