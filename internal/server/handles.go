package server

import (
	"github.com/nayakunin/gophermart/internal/storage"
)

type Server struct {
	storage *storage.DBStorage
}

func NewServer(dbStorage *storage.DBStorage) Server {
	return Server{
		storage: dbStorage,
	}
}
