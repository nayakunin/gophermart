package server

import (
	"github.com/nayakunin/gophermart/internal/config"
	"github.com/nayakunin/gophermart/internal/storage"
)

type Server struct {
	Storage *storage.DBStorage
	Cfg     config.Config
}

func NewServer(dbStorage *storage.DBStorage, cfg config.Config) Server {
	return Server{
		Storage: dbStorage,
		Cfg:     cfg,
	}
}
