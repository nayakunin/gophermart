package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStorage struct {
	Pool *pgxpool.Pool
}

func initDB(conn *pgxpool.Conn) error {
	_, err := conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
    	password VARCHAR(255) NOT NULL
	)`)
	if err != nil {
		return err
	}
	_, err = conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		status VARCHAR(255) NOT NULL,
		uploaded_at TIMESTAMP NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`)
	if err != nil {
		return err
	}
	_, err = conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS balance (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		amount FLOAT NOT NULL,
		withdrawn FLOAT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`)
	if err != nil {
		return err
	}
	_, err = conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		order_id INTEGER NOT NULL,
		amount FLOAT NOT NULL,
		processed_at TIMESTAMP NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (order_id) REFERENCES orders(id)
	)`)
	if err != nil {
		return err
	}

	return nil
}

func NewDBStorage(databaseURL string) (*DBStorage, error) {
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	err = initDB(conn)
	if err != nil {
		return nil, err
	}

	return &DBStorage{
		Pool: pool,
	}, nil
}
